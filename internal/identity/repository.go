package identity

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/AARCSX/AARCSX_Forge/internal/database"
)

// IdentityRepository defines the interface for identity persistence.
type IdentityRepository interface {
	CreateUser(ctx context.Context, u *User) error
	GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	// Session methods (for refresh tokens) can be added later
}

// IdentityRepositoryPostgres implements IdentityRepository using PostgreSQL.
type IdentityRepositoryPostgres struct {
	db *database.Postgres
}

// NewIdentityRepositoryPostgres creates a new PostgreSQL identity repository.
func NewIdentityRepositoryPostgres(db *database.Postgres) IdentityRepository {
	return &IdentityRepositoryPostgres{db: db}
}

// CreateUser inserts a new user into the database.
func (r *IdentityRepositoryPostgres) CreateUser(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, status, is_email_verified, last_login_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	_, err := r.db.DB.ExecContext(ctx, query, u.ID, u.TenantID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Status, u.IsEmailVerified, u.LastLoginAt, u.CreatedAt, u.UpdatedAt)
	return err
}

// GetUserByEmail retrieves a user by email and tenantID.
func (r *IdentityRepositoryPostgres) GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, status, is_email_verified, last_login_at, created_at, updated_at
		FROM users
		WHERE tenant_id = $1 AND email = $2
	`
	var u User
	err := r.db.DB.QueryRowContext(ctx, query, tenantID, email).Scan(
		&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Status, &u.IsEmailVerified, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// GetUserByID retrieves a user by their ID.
func (r *IdentityRepositoryPostgres) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, status, is_email_verified, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u User
	err := r.db.DB.QueryRowContext(ctx, query, userID).Scan(
		&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Status, &u.IsEmailVerified, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}