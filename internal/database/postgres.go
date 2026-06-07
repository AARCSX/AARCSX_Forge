package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/AARCSX/AARCSX_Forge/internal/config"
)

// Postgres holds the PostgreSQL connection.
type Postgres struct {
	DB *sql.DB
}

// NewPostgres creates a new PostgreSQL connection from config.
func NewPostgres(cfg config.Config) (*Postgres, error) {
	// pgx stdlib driver
	db, err := sql.Open("pgx", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("opening postgres connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("verifying postgres connection: %w", err)
	}

	return &Postgres{DB: db}, nil
}

// Close closes the database connection.
func (p *Postgres) Close() error {
	return p.DB.Close()
}