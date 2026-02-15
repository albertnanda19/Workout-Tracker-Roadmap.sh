package usecase

import (
	"context"
	"errors"
	"strings"

	"workout-tracker/internal/domain"
)

var ErrWorkoutNotFound = errors.New("workout not found")

type WorkoutUsecase struct {
	repo domain.WorkoutRepository
}

func NewWorkoutUsecase(r domain.WorkoutRepository) *WorkoutUsecase {
	return &WorkoutUsecase{repo: r}
}

func (u *WorkoutUsecase) CreatePlan(
	ctx context.Context,
	userID string,
	name string,
	notes string,
	exercises []domain.WorkoutPlanExercise,
) error {
	userID = strings.TrimSpace(userID)
	name = strings.TrimSpace(name)
	notes = strings.TrimSpace(notes)

	if userID == "" {
		return errors.New("userID is required")
	}
	if name == "" {
		return errors.New("name is required")
	}
	if len(exercises) < 1 {
		return errors.New("at least 1 exercise is required")
	}

	for _, ex := range exercises {
		if strings.TrimSpace(ex.ExerciseID) == "" {
			return errors.New("exercise_id is required")
		}
		if ex.Sets <= 0 {
			return errors.New("sets must be greater than 0")
		}
		if ex.Reps <= 0 {
			return errors.New("reps must be greater than 0")
		}
	}

	plan := &domain.WorkoutPlan{
		UserID: userID,
		Name:   name,
		Notes:  notes,
	}

	if err := u.repo.CreatePlan(ctx, plan, exercises); err != nil {
		return err
	}

	return nil
}

func (u *WorkoutUsecase) GetPlans(ctx context.Context, userID string) ([]domain.WorkoutPlan, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	return u.repo.GetPlansByUser(ctx, userID)
}

func (u *WorkoutUsecase) GetPlanByID(ctx context.Context, userID string, planID string) (*domain.WorkoutPlan, error) {
	userID = strings.TrimSpace(userID)
	planID = strings.TrimSpace(planID)

	if userID == "" {
		return nil, errors.New("userID is required")
	}
	if planID == "" {
		return nil, errors.New("planID is required")
	}

	plan, err := u.repo.GetPlanByID(ctx, planID, userID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, ErrWorkoutNotFound
	}

	return plan, nil
}

func (u *WorkoutUsecase) DeletePlan(ctx context.Context, userID string, planID string) error {
	userID = strings.TrimSpace(userID)
	planID = strings.TrimSpace(planID)

	if userID == "" {
		return errors.New("userID is required")
	}
	if planID == "" {
		return errors.New("planID is required")
	}

	if err := u.repo.DeletePlan(ctx, planID, userID); err != nil {
		return err
	}

	return nil
}
