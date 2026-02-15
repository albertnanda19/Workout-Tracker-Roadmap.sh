package usecase

import "workout-tracker/internal/domain"

type WorkoutUsecase struct {
	repo domain.WorkoutRepository
}

func NewWorkoutUsecase(r domain.WorkoutRepository) *WorkoutUsecase {
	return &WorkoutUsecase{repo: r}
}
