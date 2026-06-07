package identity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system.
type User struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Email     string    `json:"email"`
	PasswordHash string   `json:"-"` // never json'd
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Status    string    `json:"status"` // e.g., active, inactive, locked
	IsEmailVerified bool  `json:"is_email_verified"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SetPassword hashes the password and sets it on the user.
func (u *User) SetPassword(password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedBytes)
	return nil
}

// CheckPassword compares the hashed password with the provided one.
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}