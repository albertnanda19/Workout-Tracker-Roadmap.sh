package domain

import (
	"context"
	"time"
)

type WorkoutPlan struct {
	ID        string
	UserID    string
	Name      string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WorkoutPlanExercise struct {
	ID            string
	WorkoutPlanID string
	ExerciseID    string
	Sets          int
	Reps          int
	Weight        float64
	OrderIndex    int
}

type WorkoutPlanFilter struct {
	Name string
}

type WorkoutRepository interface {
	CreatePlan(ctx context.Context, plan *WorkoutPlan, exercises []WorkoutPlanExercise) error
	UpdatePlan(ctx context.Context, plan *WorkoutPlan, exercises []WorkoutPlanExercise) error
	GetPlansByUser(ctx context.Context, userID string, pagination Pagination, filters WorkoutPlanFilter) (PaginatedResult[WorkoutPlan], error)
	GetPlanByID(ctx context.Context, id string, userID string) (*WorkoutPlan, error)
	DeletePlan(ctx context.Context, id string, userID string) error
}
