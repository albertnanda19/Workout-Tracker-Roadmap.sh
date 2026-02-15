package domain

import "context"

type WorkoutPlan struct {
	ID     int64
	UserID int64
	Title  string
}

type WorkoutRepository interface {
	Create(ctx context.Context, entity *WorkoutPlan) error
}
