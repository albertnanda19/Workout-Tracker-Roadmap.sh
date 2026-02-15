package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockWorkoutPlanChecker struct {
	mock.Mock
}

func (m *MockWorkoutPlanChecker) GetOwnerID(ctx context.Context, workoutPlanID string) (string, error) {
	args := m.Called(ctx, workoutPlanID)
	return args.String(0), args.Error(1)
}
