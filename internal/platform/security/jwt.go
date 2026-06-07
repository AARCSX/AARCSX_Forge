package security

import (
	"context"
	"errors"
	"time"

	"github.com/AARCSX/AARCSX_Forge/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTKeyProvider provides signing keys for JWT tokens by key ID.
type JWTKeyProvider interface {
	// GetKey returns the key material for the given key ID.
	GetKey(keyID string) ([]byte, error)
}

// StaticKeyProvider is a simple JWTKeyProvider that uses a static set of keys.
type StaticKeyProvider struct {
	keys map[string][]byte
}

// NewStaticKeyProvider creates a new StaticKeyProvider with the given keys.
func NewStaticKeyProvider(keys map[string][]byte) *StaticKeyProvider {
	return &StaticKeyProvider{keys: keys}
}

// GetKey returns the key material for the given key ID.
func (p *StaticKeyProvider) GetKey(keyID string) ([]byte, error) {
	if key, ok := p.keys[keyID]; ok {
		return key, nil
	}
	return nil, errors.New("key not found: " + keyID)
}

// JWTTokenService implements token issuance and validation using a JWTKeyProvider.
type JWTTokenService struct {
	keyProvider JWTKeyProvider
	cfg         config.JWTConfig
}

// NewJWTTokenService creates a new JWTTokenService.
func NewJWTTokenService(keyProvider JWTKeyProvider, cfg config.JWTConfig) *JWTTokenService {
	return &JWTTokenService{
		keyProvider: keyProvider,
		cfg:         cfg,
	}
}

// IssueAccessToken creates a new access token.
// The key ID of the current signing key is placed in the "kid" header.
func (s *JWTTokenService) IssueAccessToken(ctx context.Context, subject string, claims map[string]any) (string, error) {
	// For simplicity, we assume the key provider has a single key or we use a fixed key ID.
	// In a real implementation, we would have a way to know the current key ID.
	// We'll use the key ID "current" for now and assume the key provider has it.
	currentKeyID := "current"
	key, err := s.keyProvider.GetKey(currentKeyID)
	if err != nil {
		return "", err
	}

	now := time.Now()
	tokenClaims := jwt.MapClaims{
		"iss": s.cfg.Issuer,
		"sub": subject,
		"aud": []string{s.cfg.Audience},
		"exp": jwt.NewNumericDate(now.Add(time.Duration(s.cfg.AccessTTLMin) * time.Minute)),
		"iat": jwt.NewNumericDate(now),
		"jti": uuid.New().String(),
	}

	// Add custom claims
	for k, v := range claims {
		switch k {
		case "iss", "sub", "aud", "exp", "nbf", "iat", "jti":
			// Skip reserved claims
		default:
			tokenClaims[k] = v
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	// Set the key ID in the header
	token.Header["kid"] = currentKeyID
	return token.SignedString(key)
}

// IssueRefreshToken creates a new refresh token.
func (s *JWTTokenService) IssueRefreshToken(ctx context.Context, subject string, claims map[string]any) (string, error) {
	currentKeyID := "current"
	key, err := s.keyProvider.GetKey(currentKeyID)
	if err != nil {
		return "", err
	}

	now := time.Now()
	tokenClaims := jwt.MapClaims{
		"iss": s.cfg.Issuer,
		"sub": subject,
		"aud": []string{s.cfg.Audience},
		"exp": jwt.NewNumericDate(now.Add(time.Duration(s.cfg.RefreshTTLHour) * time.Hour)),
		"iat": jwt.NewNumericDate(now),
		"jti": uuid.New().String(),
	}

	// Add custom claims
	for k, v := range claims {
		switch k {
		case "iss", "sub", "aud", "exp", "nbf", "iat", "jti":
			// Skip reserved claims
		default:
			tokenClaims[k] = v
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	token.Header["kid"] = currentKeyID
	return token.SignedString(key)
}

// ValidateToken validates a token and returns its claims.
// It extracts the key ID from the "kid" header and uses the key provider to get the key for validation.
func (s *JWTTokenService) ValidateToken(ctx context.Context, token string) (map[string]any, error) {
	// Parse the token without verifying to get the header
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract the key ID from the header
	var kid string
	if keyID, ok := parsedToken.Header["kid"].(string); ok {
		kid = keyID
	} else {
		return nil, errors.New("missing kid header in token")
	}

	// Get the key for the key ID
	key, err := s.keyProvider.GetKey(kid)
	if err != nil {
		return nil, err
	}

	// Now parse and verify the token with the key
	verifiedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !verifiedToken.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := verifiedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Convert MapClaims to map[string]any
	result := make(map[string]any)
	for k, v := range claims {
		result[k] = v
	}
	return result, nil
}