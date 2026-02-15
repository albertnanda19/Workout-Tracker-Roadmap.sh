package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"workout-tracker/internal/domain"
	"workout-tracker/internal/infrastructure/auth"
	"workout-tracker/internal/mocks"
	"workout-tracker/internal/usecase"
)

func TestUserUsecase_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputName   string
		inputEmail  string
		inputPass   string
		setupMock   func(m *mocks.MockUserRepository)
		expectedErr error
	}{
		{
			name:       "success",
			inputName:  "John",
			inputEmail: "john@example.com",
			inputPass:  "secret1",
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, sql.ErrNoRows).Once()
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
					u := args.Get(1).(*domain.User)
					assert.NotEmpty(t, u.PasswordHash)
				}).Once()
			},
			expectedErr: nil,
		},
		{
			name:       "email exists",
			inputName:  "John",
			inputEmail: "john@example.com",
			inputPass:  "secret1",
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "john@example.com").Return(&domain.User{ID: "1"}, nil).Once()
			},
			expectedErr: domain.ErrConflict,
		},
		{
			name:        "invalid input",
			inputName:   "",
			inputEmail:  "",
			inputPass:   "",
			setupMock:   func(m *mocks.MockUserRepository) {},
			expectedErr: domain.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mrepo := new(mocks.MockUserRepository)
			tt.setupMock(mrepo)

			uc := usecase.NewUserUsecase(mrepo, auth.NewJWTService("secret"))
			_, err := uc.Register(context.Background(), tt.inputName, tt.inputEmail, tt.inputPass)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
			}

			mrepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_Login_InvalidInput(t *testing.T) {
	t.Parallel()

	mrepo := new(mocks.MockUserRepository)
	uc := usecase.NewUserUsecase(mrepo, auth.NewJWTService("secret"))

	_, err := uc.Login(context.Background(), "", "pass")
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))

	_, err = uc.Login(context.Background(), "a@b.com", "")
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestUserUsecase_GetByID(t *testing.T) {
	t.Parallel()

	t.Run("invalid input", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		uc := usecase.NewUserUsecase(repo, auth.NewJWTService("secret"))
		_, err := uc.GetByID(context.Background(), "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		repo.On("GetByID", mock.Anything, "u1").Return(nil, sql.ErrNoRows).Once()

		uc := usecase.NewUserUsecase(repo, auth.NewJWTService("secret"))
		_, err := uc.GetByID(context.Background(), "u1")
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
		repo.AssertExpectations(t)
	})
}

func TestUserUsecase_Login(t *testing.T) {
	t.Parallel()

	jwtSvc := auth.NewJWTService("secret")

	tests := []struct {
		name        string
		email       string
		password    string
		setupMock   func(m *mocks.MockUserRepository)
		expectedErr error
	}{
		{
			name:     "success",
			email:    "john@example.com",
			password: "secret1",
			setupMock: func(m *mocks.MockUserRepository) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("secret1"), bcrypt.DefaultCost)
				m.On("GetByEmail", mock.Anything, "john@example.com").Return(&domain.User{ID: "u1", PasswordHash: string(hash)}, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:     "user not found",
			email:    "john@example.com",
			password: "secret1",
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, sql.ErrNoRows).Once()
			},
			expectedErr: domain.ErrUnauthorized,
		},
		{
			name:     "wrong password",
			email:    "john@example.com",
			password: "wrong",
			setupMock: func(m *mocks.MockUserRepository) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("secret1"), bcrypt.DefaultCost)
				m.On("GetByEmail", mock.Anything, "john@example.com").Return(&domain.User{ID: "u1", PasswordHash: string(hash)}, nil).Once()
			},
			expectedErr: domain.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mrepo := new(mocks.MockUserRepository)
			tt.setupMock(mrepo)

			uc := usecase.NewUserUsecase(mrepo, jwtSvc)
			tok, err := uc.Login(context.Background(), tt.email, tt.password)

			if tt.expectedErr == nil {
				require.NoError(t, err)
				require.NotNil(t, tok)
				assert.NotEmpty(t, tok.AccessToken)
			} else {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
			}

			mrepo.AssertExpectations(t)
		})
	}
}
