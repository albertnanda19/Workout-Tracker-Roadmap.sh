package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
		return fmt.Errorf("create plan: %w", domain.ErrInvalidInput)
	}
	if name == "" {
		return fmt.Errorf("create plan: %w", domain.ErrInvalidInput)
	}
	if len(exercises) < 1 {
		return fmt.Errorf("create plan: %w", domain.ErrInvalidInput)
	}

	for _, ex := range exercises {
		if strings.TrimSpace(ex.ExerciseID) == "" {
			return fmt.Errorf("create plan: %w", domain.ErrInvalidInput)
		}
		if ex.Sets <= 0 {
			return fmt.Errorf("create plan: %w", domain.ErrInvalidInput)
		}
		if ex.Reps <= 0 {
			return fmt.Errorf("create plan: %w", domain.ErrInvalidInput)
		}
	}

	plan := &domain.WorkoutPlan{
		UserID: userID,
		Name:   name,
		Notes:  notes,
	}

	if err := u.repo.CreatePlan(ctx, plan, exercises); err != nil {
		return fmt.Errorf("create plan: %w", err)
	}

	return nil
}

func (u *WorkoutUsecase) UpdatePlan(
	ctx context.Context,
	userID string,
	planID string,
	name string,
	notes string,
	exercises []domain.WorkoutPlanExercise,
) error {
	userID = strings.TrimSpace(userID)
	planID = strings.TrimSpace(planID)
	name = strings.TrimSpace(name)
	notes = strings.TrimSpace(notes)

	if userID == "" {
		return fmt.Errorf("update plan: %w", domain.ErrInvalidInput)
	}
	if planID == "" {
		return fmt.Errorf("update plan: %w", domain.ErrInvalidInput)
	}
	if name == "" {
		return fmt.Errorf("update plan: %w", domain.ErrInvalidInput)
	}
	if len(exercises) < 1 {
		return fmt.Errorf("update plan: %w", domain.ErrInvalidInput)
	}

	for _, ex := range exercises {
		if strings.TrimSpace(ex.ExerciseID) == "" {
			return fmt.Errorf("update plan: %w", domain.ErrInvalidInput)
		}
		if ex.Sets <= 0 {
			return fmt.Errorf("update plan: %w", domain.ErrInvalidInput)
		}
		if ex.Reps <= 0 {
			return fmt.Errorf("update plan: %w", domain.ErrInvalidInput)
		}
	}

	plan := &domain.WorkoutPlan{
		ID:     planID,
		UserID: userID,
		Name:   name,
		Notes:  notes,
	}

	if err := u.repo.UpdatePlan(ctx, plan, exercises); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("update plan: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("update plan: %w", err)
	}

	return nil
}

func (u *WorkoutUsecase) GetPlans(ctx context.Context, userID string) ([]domain.WorkoutPlan, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("get plans: %w", domain.ErrInvalidInput)
	}

	plans, err := u.repo.GetPlansByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get plans: %w", err)
	}
	return plans, nil
}

func (u *WorkoutUsecase) GetPlanByID(ctx context.Context, userID string, planID string) (*domain.WorkoutPlan, error) {
	userID = strings.TrimSpace(userID)
	planID = strings.TrimSpace(planID)

	if userID == "" {
		return nil, fmt.Errorf("get plan: %w", domain.ErrInvalidInput)
	}
	if planID == "" {
		return nil, fmt.Errorf("get plan: %w", domain.ErrInvalidInput)
	}

	plan, err := u.repo.GetPlanByID(ctx, planID, userID)
	if err != nil {
		return nil, fmt.Errorf("get plan: %w", err)
	}
	if plan == nil {
		return nil, fmt.Errorf("get plan: %w", domain.ErrNotFound)
	}

	return plan, nil
}

func (u *WorkoutUsecase) DeletePlan(ctx context.Context, userID string, planID string) error {
	userID = strings.TrimSpace(userID)
	planID = strings.TrimSpace(planID)

	if userID == "" {
		return fmt.Errorf("delete plan: %w", domain.ErrInvalidInput)
	}
	if planID == "" {
		return fmt.Errorf("delete plan: %w", domain.ErrInvalidInput)
	}

	if err := u.repo.DeletePlan(ctx, planID, userID); err != nil {
		return fmt.Errorf("delete plan: %w", err)
	}

	return nil
}
