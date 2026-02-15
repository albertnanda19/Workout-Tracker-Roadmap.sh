package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"workout-tracker/internal/domain"
	"workout-tracker/internal/infrastructure/auth"
)

type UserUsecase struct {
	repo domain.UserRepository
	jwt  *auth.JWTService
}

func NewUserUsecase(r domain.UserRepository, jwt *auth.JWTService) *UserUsecase {
	return &UserUsecase{repo: r, jwt: jwt}
}

func (u *UserUsecase) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return nil, fmt.Errorf("register: %w", domain.ErrInvalidInput)
	}
	if email == "" {
		return nil, fmt.Errorf("register: %w", domain.ErrInvalidInput)
	}
	if len(password) < 6 {
		return nil, fmt.Errorf("register: %w", domain.ErrInvalidInput)
	}

	_, err := u.repo.GetByEmail(ctx, email)
	if err == nil {
		return nil, fmt.Errorf("register: %w", domain.ErrConflict)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("register: %w", err)
	}

	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(h),
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	user.PasswordHash = ""
	return user, nil
}

func (u *UserUsecase) Login(ctx context.Context, email, password string) (*domain.AuthToken, error) {
	email = strings.TrimSpace(email)

	if email == "" {
		return nil, fmt.Errorf("login: %w", domain.ErrInvalidInput)
	}
	if password == "" {
		return nil, fmt.Errorf("login: %w", domain.ErrInvalidInput)
	}

	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("login: %w", domain.ErrUnauthorized)
		}
		return nil, fmt.Errorf("login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("login: %w", domain.ErrUnauthorized)
	}

	token, exp, err := u.jwt.Generate(user.ID)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	return &domain.AuthToken{AccessToken: token, ExpiresAt: exp}, nil
}

func (u *UserUsecase) GetByID(ctx context.Context, id string) (*domain.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("get user: %w", domain.ErrInvalidInput)
	}

	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get user: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	user.PasswordHash = ""
	return user, nil
}
