package identity

import "context"

type Service interface {
	RegisterUser(ctx context.Context, in RegisterUserInput) (User, error)
	Authenticate(ctx context.Context, in AuthenticateInput) (TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (TokenPair, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user User) (User, error)
	GetUserByEmail(ctx context.Context, tenantID, email string) (User, error)
	StoreSession(ctx context.Context, s Session) error
	GetSession(ctx context.Context, tokenID string) (Session, error)
}

type PermissionChecker interface {
	Can(ctx context.Context, actorID, tenantID, action string) (bool, error)
}

type RegisterUserInput struct {
	TenantID string
	Email    string
	Password string
}

type AuthenticateInput struct {
	TenantID string
	Email    string
	Password string
}

type User struct {
	ID       string
	TenantID string
	Email    string
}

type Session struct {
	TokenID  string
	UserID   string
	TenantID string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
