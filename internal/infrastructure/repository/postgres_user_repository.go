package repository

import (
	"context"
	"database/sql"
	"fmt"

	"workout-tracker/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) domain.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	if user == nil {
		return fmt.Errorf("create user: user is nil")
	}

	const q = `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
	`

	if _, err := r.db.ExecContext(ctx, q, user.Name, user.Email, user.PasswordHash); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var u domain.User
	if err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &u, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	const q = `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var u domain.User
	if err := r.db.QueryRowContext(ctx, q, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &u, nil
}
