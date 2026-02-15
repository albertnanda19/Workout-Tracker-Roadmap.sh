package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"workout-tracker/internal/domain"
	"workout-tracker/internal/mocks"
	"workout-tracker/internal/usecase"
)

func TestWorkoutUsecase_CreatePlan(t *testing.T) {
	t.Parallel()

	validExercises := []domain.WorkoutPlanExercise{{ExerciseID: "e1", Sets: 3, Reps: 10}}

	tests := []struct {
		name        string
		userID      string
		planName    string
		exercises   []domain.WorkoutPlanExercise
		setupMock   func(m *mocks.MockWorkoutRepository)
		expectedErr error
	}{
		{
			name:      "success",
			userID:    "u1",
			planName:  "Plan",
			exercises: validExercises,
			setupMock: func(m *mocks.MockWorkoutRepository) {
				m.On("CreatePlan", mock.Anything, mock.AnythingOfType("*domain.WorkoutPlan"), validExercises).Return(nil).Once()
			},
		},
		{
			name:        "name empty",
			userID:      "u1",
			planName:    "",
			exercises:   validExercises,
			setupMock:   func(m *mocks.MockWorkoutRepository) {},
			expectedErr: domain.ErrInvalidInput,
		},
		{
			name:        "exercise empty",
			userID:      "u1",
			planName:    "Plan",
			exercises:   nil,
			setupMock:   func(m *mocks.MockWorkoutRepository) {},
			expectedErr: domain.ErrInvalidInput,
		},
		{
			name:        "sets invalid",
			userID:      "u1",
			planName:    "Plan",
			exercises:   []domain.WorkoutPlanExercise{{ExerciseID: "e1", Sets: 0, Reps: 10}},
			setupMock:   func(m *mocks.MockWorkoutRepository) {},
			expectedErr: domain.ErrInvalidInput,
		},
		{
			name:        "reps invalid",
			userID:      "u1",
			planName:    "Plan",
			exercises:   []domain.WorkoutPlanExercise{{ExerciseID: "e1", Sets: 1, Reps: 0}},
			setupMock:   func(m *mocks.MockWorkoutRepository) {},
			expectedErr: domain.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.MockWorkoutRepository)
			tt.setupMock(repo)

			uc := usecase.NewWorkoutUsecase(repo)
			err := uc.CreatePlan(context.Background(), tt.userID, tt.planName, "", tt.exercises)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestWorkoutUsecase_UpdatePlan_InvalidInput(t *testing.T) {
	t.Parallel()

	repo := new(mocks.MockWorkoutRepository)
	uc := usecase.NewWorkoutUsecase(repo)

	bad := []struct {
		name string
		err  error
	}{
		{"missing user", uc.UpdatePlan(context.Background(), "", "p1", "name", "", []domain.WorkoutPlanExercise{{ExerciseID: "e1", Sets: 1, Reps: 1}})},
		{"missing plan", uc.UpdatePlan(context.Background(), "u1", "", "name", "", []domain.WorkoutPlanExercise{{ExerciseID: "e1", Sets: 1, Reps: 1}})},
		{"missing name", uc.UpdatePlan(context.Background(), "u1", "p1", "", "", []domain.WorkoutPlanExercise{{ExerciseID: "e1", Sets: 1, Reps: 1}})},
		{"no exercises", uc.UpdatePlan(context.Background(), "u1", "p1", "name", "", nil)},
		{"bad sets", uc.UpdatePlan(context.Background(), "u1", "p1", "name", "", []domain.WorkoutPlanExercise{{ExerciseID: "e1", Sets: 0, Reps: 1}})},
		{"bad reps", uc.UpdatePlan(context.Background(), "u1", "p1", "name", "", []domain.WorkoutPlanExercise{{ExerciseID: "e1", Sets: 1, Reps: 0}})},
	}

	for _, tt := range bad {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, tt.err)
			assert.True(t, errors.Is(tt.err, domain.ErrInvalidInput))
		})
	}
}

func TestWorkoutUsecase_UpdatePlan(t *testing.T) {
	t.Parallel()

	exercises := []domain.WorkoutPlanExercise{{ExerciseID: "e1", Sets: 3, Reps: 10}}

	tests := []struct {
		name        string
		setupMock   func(m *mocks.MockWorkoutRepository)
		expectedErr error
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockWorkoutRepository) {
				m.On("UpdatePlan", mock.Anything, mock.AnythingOfType("*domain.WorkoutPlan"), exercises).Return(nil).Once()
			},
		},
		{
			name: "not found",
			setupMock: func(m *mocks.MockWorkoutRepository) {
				m.On("UpdatePlan", mock.Anything, mock.AnythingOfType("*domain.WorkoutPlan"), exercises).Return(sql.ErrNoRows).Once()
			},
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.MockWorkoutRepository)
			tt.setupMock(repo)

			uc := usecase.NewWorkoutUsecase(repo)
			err := uc.UpdatePlan(context.Background(), "u1", "p1", "name", "", exercises)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestWorkoutUsecase_DeletePlan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupMock   func(m *mocks.MockWorkoutRepository)
		expectedErr error
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockWorkoutRepository) {
				m.On("DeletePlan", mock.Anything, "p1", "u1").Return(nil).Once()
			},
		},
		{
			name: "not found",
			setupMock: func(m *mocks.MockWorkoutRepository) {
				m.On("DeletePlan", mock.Anything, "p1", "u1").Return(sql.ErrNoRows).Once()
			},
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.MockWorkoutRepository)
			tt.setupMock(repo)

			uc := usecase.NewWorkoutUsecase(repo)
			err := uc.DeletePlan(context.Background(), "u1", "p1")
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestWorkoutUsecase_GetPlans(t *testing.T) {
	t.Parallel()

	repo := new(mocks.MockWorkoutRepository)
	uc := usecase.NewWorkoutUsecase(repo)

	expected := []domain.WorkoutPlan{{ID: "p1"}}
	repo.On("GetPlansByUser", mock.Anything, "u1").Return(expected, nil).Once()

	plans, err := uc.GetPlans(context.Background(), "u1")
	require.NoError(t, err)
	assert.Len(t, plans, 1)
	repo.AssertExpectations(t)

	repo2 := new(mocks.MockWorkoutRepository)
	uc2 := usecase.NewWorkoutUsecase(repo2)
	repo2.On("GetPlansByUser", mock.Anything, "u1").Return(nil, errors.New("db"))

	_, err = uc2.GetPlans(context.Background(), "u1")
	require.Error(t, err)
}
