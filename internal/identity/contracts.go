package identity

import (
	"context"

	"github.com/google/uuid"
	"errors"
)

// Errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// Service defines the contract for identity business logic.
type Service interface {
	RegisterUser(ctx context.Context, in *RegisterUserInput) (*UserResponse, error)
	Authenticate(ctx context.Context, in *AuthenticateInput) (*TokenPairResponse, error)
	RefreshToken(ctx context.Context, in *RefreshTokenInput) (*TokenPairResponse, error)
}

// Repository defines the contract for identity persistence.
type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
}

// PermissionChecker defines the contract for checking permissions.
type PermissionChecker interface {
	Can(ctx context.Context, actorID, tenantID, action string) (bool, error)
}