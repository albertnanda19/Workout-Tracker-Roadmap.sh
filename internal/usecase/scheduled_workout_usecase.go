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
	repo repository.ScheduledWorkoutRepository
	db   *sql.DB
}

func NewScheduledWorkoutUsecase(repo repository.ScheduledWorkoutRepository, db *sql.DB) *ScheduledWorkoutUsecase {
	return &ScheduledWorkoutUsecase{repo: repo, db: db}
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

	if u.db != nil {
		var ok bool
		if err := u.db.QueryRowContext(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM workout_plans WHERE id = $1 AND user_id = $2
			)
		`, workoutPlanID, userID).Scan(&ok); err != nil {
			return fmt.Errorf("schedule workout: %w", err)
		}
		if !ok {
			return fmt.Errorf("schedule workout: %w", domain.ErrForbidden)
		}
	}

	existing, err := u.repo.GetByUserAndDate(ctx, userID, date)
	if err != nil {
		return fmt.Errorf("schedule workout: %w", err)
	}
	for _, sw := range existing {
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

func (u *ScheduledWorkoutUsecase) GetUserScheduleByDate(ctx context.Context, userID string, date time.Time) ([]domain.ScheduledWorkout, error) {
	if userID == "" {
		return nil, fmt.Errorf("get schedules: %w", domain.ErrInvalidInput)
	}

	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	items, err := u.repo.GetByUserAndDate(ctx, userID, d)
	if err != nil {
		return nil, fmt.Errorf("get schedules: %w", err)
	}
	return items, nil
}

func (u *ScheduledWorkoutUsecase) GetAllUserSchedules(ctx context.Context, userID string) ([]domain.ScheduledWorkout, error) {
	if userID == "" {
		return nil, fmt.Errorf("get schedules: %w", domain.ErrInvalidInput)
	}

	items, err := u.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get schedules: %w", err)
	}
	return items, nil
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
