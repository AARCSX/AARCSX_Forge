package security

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher implements PasswordHasher using bcrypt.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new bcrypt hasher with the given cost.
// If cost is less than bcrypt.MinCost, it will be set to bcrypt.DefaultCost.
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

// Hash generates a bcrypt hash of the password.
func (h *BcryptHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Compare checks if the password matches the hash.
func (h *BcryptHasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}