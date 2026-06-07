package identity

// RegisterUserInput represents the input for user registration.
type RegisterUserInput struct {
	TenantID string `json:"tenant_id" validate:"required,uuid"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// AuthenticateInput represents the input for user authentication.
type AuthenticateInput struct {
	TenantID string `json:"tenant_id" validate:"required,uuid"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenInput represents the input for refreshing a token.
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UserResponse represents the user in responses.
type UserResponse struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
}

// TokenPairResponse represents the token pair in responses.
type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}