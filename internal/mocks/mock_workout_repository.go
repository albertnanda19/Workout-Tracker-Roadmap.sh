package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"workout-tracker/internal/domain"
)

type MockWorkoutRepository struct {
	mock.Mock
}

func (m *MockWorkoutRepository) CreatePlan(ctx context.Context, plan *domain.WorkoutPlan, exercises []domain.WorkoutPlanExercise) error {
	args := m.Called(ctx, plan, exercises)
	return args.Error(0)
}

func (m *MockWorkoutRepository) UpdatePlan(ctx context.Context, plan *domain.WorkoutPlan, exercises []domain.WorkoutPlanExercise) error {
	args := m.Called(ctx, plan, exercises)
	return args.Error(0)
}

func (m *MockWorkoutRepository) GetPlansByUser(ctx context.Context, userID string) ([]domain.WorkoutPlan, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.WorkoutPlan), args.Error(1)
}

func (m *MockWorkoutRepository) GetPlanByID(ctx context.Context, id string, userID string) (*domain.WorkoutPlan, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WorkoutPlan), args.Error(1)
}

func (m *MockWorkoutRepository) DeletePlan(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}
