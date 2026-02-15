package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"workout-tracker/internal/domain"
	"workout-tracker/internal/repository"
)

var (
	ErrScheduleConflict = errors.New("schedule already exists")
	ErrScheduleNotFound = errors.New("schedule not found")
	ErrForbiddenPlan    = errors.New("forbidden")
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
		return errors.New("userID is required")
	}
	if workoutPlanID == "" {
		return errors.New("workoutPlanID is required")
	}

	date := time.Date(scheduledDate.Year(), scheduledDate.Month(), scheduledDate.Day(), 0, 0, 0, 0, time.UTC)
	today := time.Now().UTC()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	if date.Before(today) {
		return errors.New("scheduled date must be today or later")
	}

	if u.db != nil {
		var ok bool
		if err := u.db.QueryRowContext(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM workout_plans WHERE id = $1 AND user_id = $2
			)
		`, workoutPlanID, userID).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return ErrForbiddenPlan
		}
	}

	existing, err := u.repo.GetByUserAndDate(ctx, userID, date)
	if err != nil {
		return err
	}
	for _, sw := range existing {
		if sw.WorkoutPlanID == workoutPlanID {
			return ErrScheduleConflict
		}
	}

	sw := &domain.ScheduledWorkout{
		UserID:        userID,
		WorkoutPlanID: workoutPlanID,
		ScheduledDate: date,
	}

	if err := u.repo.Create(ctx, sw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrScheduleConflict
		}
		return err
	}

	return nil
}

func (u *ScheduledWorkoutUsecase) GetUserScheduleByDate(ctx context.Context, userID string, date time.Time) ([]domain.ScheduledWorkout, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	return u.repo.GetByUserAndDate(ctx, userID, d)
}

func (u *ScheduledWorkoutUsecase) GetAllUserSchedules(ctx context.Context, userID string) ([]domain.ScheduledWorkout, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	return u.repo.GetByUser(ctx, userID)
}

func (u *ScheduledWorkoutUsecase) DeleteSchedule(ctx context.Context, id, userID string) error {
	if userID == "" {
		return errors.New("userID is required")
	}
	if id == "" {
		return errors.New("id is required")
	}

	if err := u.repo.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrScheduleNotFound
		}
		return err
	}

	return nil
}
