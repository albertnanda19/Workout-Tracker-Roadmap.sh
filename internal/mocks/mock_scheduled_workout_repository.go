package mocks

import (
	"context"

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

func (m *MockScheduledWorkoutRepository) GetByUser(ctx context.Context, userID string, pagination domain.Pagination, filters domain.ScheduledWorkoutFilter) (domain.PaginatedResult[domain.ScheduledWorkout], error) {
	args := m.Called(ctx, userID, pagination, filters)
	if args.Get(0) == nil {
		return domain.PaginatedResult[domain.ScheduledWorkout]{}, args.Error(1)
	}
	return args.Get(0).(domain.PaginatedResult[domain.ScheduledWorkout]), args.Error(1)
}

func (m *MockScheduledWorkoutRepository) Delete(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}
