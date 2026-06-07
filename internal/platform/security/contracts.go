package security

import "context"

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenService interface {
	IssueAccessToken(ctx context.Context, subject string, claims map[string]any) (string, error)
	IssueRefreshToken(ctx context.Context, subject string, claims map[string]any) (string, error)
	ValidateToken(ctx context.Context, token string) (map[string]any, error)
}
