package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"workout-tracker/internal/domain"
	"workout-tracker/internal/mocks"
	"workout-tracker/internal/usecase"
)

func TestWorkoutUsecase_GetPlanByID(t *testing.T) {
	t.Parallel()

	t.Run("invalid input", func(t *testing.T) {
		repo := new(mocks.MockWorkoutRepository)
		uc := usecase.NewWorkoutUsecase(repo)
		_, err := uc.GetPlanByID(context.Background(), "", "p1")
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrInvalidInput))
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mocks.MockWorkoutRepository)
		repo.On("GetPlanByID", mock.Anything, "p1", "u1").Return(nil, nil).Once()
		uc := usecase.NewWorkoutUsecase(repo)
		_, err := uc.GetPlanByID(context.Background(), "u1", "p1")
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrNotFound))
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockWorkoutRepository)
		repo.On("GetPlanByID", mock.Anything, "p1", "u1").Return(&domain.WorkoutPlan{ID: "p1"}, nil).Once()
		uc := usecase.NewWorkoutUsecase(repo)
		p, err := uc.GetPlanByID(context.Background(), "u1", "p1")
		require.NoError(t, err)
		require.NotNil(t, p)
		assert.Equal(t, "p1", p.ID)
		repo.AssertExpectations(t)
	})
}

func TestWorkoutUsecase_DeletePlan_InvalidInput(t *testing.T) {
	t.Parallel()

	repo := new(mocks.MockWorkoutRepository)
	uc := usecase.NewWorkoutUsecase(repo)

	err := uc.DeletePlan(context.Background(), "", "p1")
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))

	err = uc.DeletePlan(context.Background(), "u1", "")
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}
