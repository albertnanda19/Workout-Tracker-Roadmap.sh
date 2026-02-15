package repository

import (
	"context"

	"workout-tracker/internal/domain"
)

type ScheduledWorkoutRepository interface {
	Create(ctx context.Context, sw *domain.ScheduledWorkout) error
	GetByUser(ctx context.Context, userID string, pagination domain.Pagination, filters domain.ScheduledWorkoutFilter) (domain.PaginatedResult[domain.ScheduledWorkout], error)
	Delete(ctx context.Context, id string, userID string) error
}

type WorkoutPlanChecker interface {
	GetOwnerID(ctx context.Context, workoutPlanID string) (string, error)
}
