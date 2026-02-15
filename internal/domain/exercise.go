package domain

import "context"

type Exercise struct {
	ID          string
	Name        string
	Description string
	Category    string
	MuscleGroup string
}

type ExerciseRepository interface {
	GetAll(ctx context.Context) ([]Exercise, error)
	GetByID(ctx context.Context, id string) (*Exercise, error)
}
