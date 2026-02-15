package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"workout-tracker/internal/domain"
	"workout-tracker/internal/repository"
)

type ScheduledWorkoutUsecase struct {
	repo        repository.ScheduledWorkoutRepository
	planChecker repository.WorkoutPlanChecker
}

func NewScheduledWorkoutUsecase(repo repository.ScheduledWorkoutRepository, planChecker repository.WorkoutPlanChecker) *ScheduledWorkoutUsecase {
	return &ScheduledWorkoutUsecase{repo: repo, planChecker: planChecker}
}

func (u *ScheduledWorkoutUsecase) ScheduleWorkout(ctx context.Context, userID, workoutPlanID string, scheduledDate time.Time) error {
	if userID == "" {
		return fmt.Errorf("schedule workout: %w", domain.ErrInvalidInput)
	}
	if workoutPlanID == "" {
		return fmt.Errorf("schedule workout: %w", domain.ErrInvalidInput)
	}

	date := time.Date(scheduledDate.Year(), scheduledDate.Month(), scheduledDate.Day(), 0, 0, 0, 0, time.UTC)
	today := time.Now().UTC()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	if date.Before(today) {
		return fmt.Errorf("schedule workout: %w", domain.ErrInvalidInput)
	}

	ownerID, err := u.planChecker.GetOwnerID(ctx, workoutPlanID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("schedule workout: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("schedule workout: %w", err)
	}
	if ownerID != userID {
		return fmt.Errorf("schedule workout: %w", domain.ErrForbidden)
	}

	res, err := u.repo.GetByUser(ctx, userID, domain.NewPagination(1, 100), domain.ScheduledWorkoutFilter{Date: &date})
	if err != nil {
		return fmt.Errorf("schedule workout: %w", err)
	}
	for _, sw := range res.Data {
		if sw.WorkoutPlanID == workoutPlanID {
			return fmt.Errorf("schedule workout: %w", domain.ErrConflict)
		}
	}

	sw := &domain.ScheduledWorkout{
		UserID:        userID,
		WorkoutPlanID: workoutPlanID,
		ScheduledDate: date,
	}

	if err := u.repo.Create(ctx, sw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("schedule workout: %w", domain.ErrConflict)
		}
		return fmt.Errorf("schedule workout: %w", err)
	}

	return nil
}

func (u *ScheduledWorkoutUsecase) GetSchedules(ctx context.Context, userID string, pagination domain.Pagination, filters domain.ScheduledWorkoutFilter) (domain.PaginatedResult[domain.ScheduledWorkout], error) {
	if userID == "" {
		return domain.PaginatedResult[domain.ScheduledWorkout]{}, fmt.Errorf("get schedules: %w", domain.ErrInvalidInput)
	}

	res, err := u.repo.GetByUser(ctx, userID, pagination, filters)
	if err != nil {
		return domain.PaginatedResult[domain.ScheduledWorkout]{}, fmt.Errorf("get schedules: %w", err)
	}
	return res, nil
}

func (u *ScheduledWorkoutUsecase) DeleteSchedule(ctx context.Context, id, userID string) error {
	if userID == "" {
		return fmt.Errorf("delete schedule: %w", domain.ErrInvalidInput)
	}
	if id == "" {
		return fmt.Errorf("delete schedule: %w", domain.ErrInvalidInput)
	}

	if err := u.repo.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("delete schedule: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("delete schedule: %w", err)
	}

	return nil
}
