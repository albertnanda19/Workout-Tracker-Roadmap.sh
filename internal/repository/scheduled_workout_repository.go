package repository

import (
	"context"
	"time"

	"workout-tracker/internal/domain"
)

type ScheduledWorkoutRepository interface {
	Create(ctx context.Context, sw *domain.ScheduledWorkout) error
	GetByUserAndDate(ctx context.Context, userID string, date time.Time) ([]domain.ScheduledWorkout, error)
	GetByUser(ctx context.Context, userID string) ([]domain.ScheduledWorkout, error)
	Delete(ctx context.Context, id string, userID string) error
}

type WorkoutPlanChecker interface {
	GetOwnerID(ctx context.Context, workoutPlanID string) (string, error)
}
