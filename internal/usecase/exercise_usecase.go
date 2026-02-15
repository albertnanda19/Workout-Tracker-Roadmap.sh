package usecase

import (
	"context"
	"errors"

	"workout-tracker/internal/domain"
)

type ExerciseUsecase struct {
	repo domain.ExerciseRepository
}

func NewExerciseUsecase(r domain.ExerciseRepository) *ExerciseUsecase {
	return &ExerciseUsecase{repo: r}
}

func (u *ExerciseUsecase) GetAll(ctx context.Context) ([]domain.Exercise, error) {
	if u == nil {
		return nil, errors.New("exercise usecase is nil")
	}
	return u.repo.GetAll(ctx)
}
