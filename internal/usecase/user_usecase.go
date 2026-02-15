package usecase

import (
	"context"
	"database/sql"
	"errors"
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
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	_, err := u.repo.GetByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrEmailAlreadyExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(h),
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	return user, nil
}

func (u *UserUsecase) Login(ctx context.Context, email, password string) (*domain.AuthToken, error) {
	email = strings.TrimSpace(email)

	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, exp, err := u.jwt.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	return &domain.AuthToken{AccessToken: token, ExpiresAt: exp}, nil
}
