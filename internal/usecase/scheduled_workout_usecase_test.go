package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"workout-tracker/internal/domain"
	"workout-tracker/internal/mocks"
	"workout-tracker/internal/usecase"
)

func TestScheduledWorkoutUsecase_ScheduleWorkout(t *testing.T) {
	t.Parallel()

	today := time.Now().UTC()
	tomorrow := time.Date(today.Year(), today.Month(), today.Day()+1, 0, 0, 0, 0, time.UTC)
	yesterday := time.Date(today.Year(), today.Month(), today.Day()-1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		date        time.Time
		ownerID     string
		ownerErr    error
		existing    []domain.ScheduledWorkout
		existingErr error
		createErr   error
		expectedErr error
	}{
		{
			name:    "success",
			date:    tomorrow,
			ownerID: "u1",
		},
		{
			name:        "past date",
			date:        yesterday,
			ownerID:     "u1",
			expectedErr: domain.ErrInvalidInput,
		},
		{
			name:        "plan not found",
			date:        tomorrow,
			ownerErr:    sql.ErrNoRows,
			expectedErr: domain.ErrNotFound,
		},
		{
			name:        "unauthorized plan",
			date:        tomorrow,
			ownerID:     "someone-else",
			expectedErr: domain.ErrForbidden,
		},
		{
			name:        "duplicate schedule",
			date:        tomorrow,
			ownerID:     "u1",
			existing:    []domain.ScheduledWorkout{{WorkoutPlanID: "p1"}},
			expectedErr: domain.ErrConflict,
		},
		{
			name:        "duplicate from repo create",
			date:        tomorrow,
			ownerID:     "u1",
			createErr:   sql.ErrNoRows,
			expectedErr: domain.ErrConflict,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.MockScheduledWorkoutRepository)
			checker := new(mocks.MockWorkoutPlanChecker)

			if !errors.Is(tt.expectedErr, domain.ErrInvalidInput) {
				checker.On("GetOwnerID", mock.Anything, "p1").Return(tt.ownerID, tt.ownerErr).Once()
			}

			if tt.expectedErr == nil || errors.Is(tt.expectedErr, domain.ErrConflict) {
				repo.On("GetByUser", mock.Anything, "u1", mock.Anything, mock.Anything).
					Return(domain.NewPaginatedResult(tt.existing, len(tt.existing), domain.NewPagination(1, 100)), tt.existingErr)
			}

			if tt.expectedErr == nil || tt.createErr != nil {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.ScheduledWorkout")).Return(tt.createErr)
			}

			uc := usecase.NewScheduledWorkoutUsecase(repo, checker)
			err := uc.ScheduleWorkout(context.Background(), "u1", "p1", tt.date)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
			}

			repo.AssertExpectations(t)
			checker.AssertExpectations(t)
		})
	}
}

func TestScheduledWorkoutUsecase_DeleteSchedule(t *testing.T) {
	t.Parallel()

	repo := new(mocks.MockScheduledWorkoutRepository)
	checker := new(mocks.MockWorkoutPlanChecker)
	uc := usecase.NewScheduledWorkoutUsecase(repo, checker)

	repo.On("Delete", mock.Anything, "s1", "u1").Return(nil).Once()
	err := uc.DeleteSchedule(context.Background(), "s1", "u1")
	require.NoError(t, err)
	repo.AssertExpectations(t)

	repo2 := new(mocks.MockScheduledWorkoutRepository)
	uc2 := usecase.NewScheduledWorkoutUsecase(repo2, checker)
	repo2.On("Delete", mock.Anything, "s1", "u1").Return(sql.ErrNoRows).Once()
	err = uc2.DeleteSchedule(context.Background(), "s1", "u1")
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}
