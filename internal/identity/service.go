package identity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/AARCSX/AARCSX_Forge/internal/config"
	"github.com/AARCSX/AARCSX_Forge/internal/database"
	"github.com/AARCSX/AARCSX_Forge/internal/logger"
	"github.com/AARCSX/AARCSX_Forge/internal/platform/security"
)

// IdentityService provides identity business logic.
type IdentityService struct {
	repo        		Repository
	cfg            	config.Config
	logger         	*logger.Logger
	permissionChecker PermissionChecker
	db            	*database.Postgres
	tokenService    *security.JWTTokenService
}

// NewIdentityService creates a new identity service.
func NewIdentityService(repo Repository, cfg config.Config, logger *logger.Logger, permissionChecker PermissionChecker, db *database.Postgres) *IdentityService {
	// Create a simple static key provider for JWT signing
	// In production, this should be loaded from secure storage or a key management system
	staticKeys := map[string][]byte{
		"current": []byte(cfg.JWT.SigningKey),
	}
	keyProvider := security.NewStaticKeyProvider(staticKeys)
	tokenService := security.NewJWTTokenService(keyProvider, cfg.JWT)

	return &IdentityService{
		repo:             repo,
		cfg:              cfg,
		logger:           logger,
		permissionChecker: permissionChecker,
		db:               db,
		tokenService:     tokenService,
	}
}

// RegisterUser creates a new user.
func (s *IdentityService) RegisterUser(ctx context.Context, in *RegisterUserInput) (*UserResponse, error) {
	// Check if user already exists
	tenantID, err := uuid.Parse(in.TenantID)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.GetUserByEmail(ctx, tenantID, in.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	// Create new user
	user := &User{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Email:        in.Email,
		FirstName:    "", // TODO: get from input if needed
		LastName:     "",
		Status:       "active",
		IsEmailVerified: false,
	}

	// Hash password
	if err := user.SetPassword(in.Password); err != nil {
		return nil, err
	}

	// Persist user
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID.String(),
		TenantID: user.TenantID.String(),
		Email:    user.Email,
	}, nil
}

// Authenticate authenticates a user and returns a token pair.
func (s *IdentityService) Authenticate(ctx context.Context, in *AuthenticateInput) (*TokenPairResponse, error) {
	tenantID, err := uuid.Parse(in.TenantID)
	if err != nil {
		return nil, err
	}

	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, tenantID, in.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if err := user.CheckPassword(in.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Store refresh token hash in sessions table for validation/revocation
	hashedToken := hashToken(refreshToken)
	_, err = s.db.DB.ExecContext(ctx, `
		INSERT INTO sessions (id, user_id, refresh_token_hash, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, uuid.New(), user.ID, hashedToken, "", "", time.Now().Add(time.Duration(s.cfg.JWT.RefreshTTLHour) * time.Hour))
	if err != nil {
		return nil, err
	}

	return &TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token.
func (s *IdentityService) RefreshToken(ctx context.Context, in *RefreshTokenInput) (*TokenPairResponse, error) {
	// Parse and validate the refresh token
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(in.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.SigningKey), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidCredentials
	}

	// Extract user ID from claims
	userIDStr := claims.Subject // Subject is a string
	if userIDStr == "" {
		return nil, ErrInvalidCredentials
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Validate refresh token hasn't been revoked (check sessions table)
	hashedToken := hashToken(in.RefreshToken)
	var exists bool
	err = s.db.DB.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM sessions
			WHERE user_id = $1
			  AND refresh_token_hash = $2
			  AND revoked_at IS NULL
			  AND expires_at > NOW()
		)
	`, userID, hashedToken).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrInvalidCredentials
	}

	// Get user with tenant ID
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Generate new tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateAccessToken creates a JWT access token.
func (s *IdentityService) generateAccessToken(u *User) (string, error) {
	accessClaims := map[string]any{
		"role": strings.Join(s.getUserRoles(u.ID.String()), ","),
	}
	accessToken, err := s.tokenService.IssueAccessToken(context.Background(), u.ID.String(), accessClaims)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// generateRefreshToken creates a JWT refresh token.
func (s *IdentityService) generateRefreshToken(u *User) (string, error) {
	refreshClaims := map[string]any{
		"role": strings.Join(s.getUserRoles(u.ID.String()), ","),
	}
	refreshToken, err := s.tokenService.IssueRefreshToken(context.Background(), u.ID.String(), refreshClaims)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

// hashToken hashes a token using SHA-256 for secure storage.
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// getUserRoles returns role identifiers for a given user ID.
// Placeholder implementation; replace with actual role lookup.
func (s *IdentityService) getUserRoles(_ string) []string {
    // TODO: integrate with permissionChecker or role service.
    return []string{}
}