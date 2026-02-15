package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"workout-tracker/internal/domain"
)

type MockScheduledWorkoutRepository struct {
	mock.Mock
}

func (m *MockScheduledWorkoutRepository) Create(ctx context.Context, sw *domain.ScheduledWorkout) error {
	args := m.Called(ctx, sw)
	return args.Error(0)
}

func (m *MockScheduledWorkoutRepository) GetByUserAndDate(ctx context.Context, userID string, date time.Time) ([]domain.ScheduledWorkout, error) {
	args := m.Called(ctx, userID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ScheduledWorkout), args.Error(1)
}

func (m *MockScheduledWorkoutRepository) GetByUser(ctx context.Context, userID string) ([]domain.ScheduledWorkout, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ScheduledWorkout), args.Error(1)
}

func (m *MockScheduledWorkoutRepository) Delete(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}
