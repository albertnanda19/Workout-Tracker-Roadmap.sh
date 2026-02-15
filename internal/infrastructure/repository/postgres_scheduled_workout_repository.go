package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"workout-tracker/internal/domain"
	irepo "workout-tracker/internal/repository"
)

type PostgresScheduledWorkoutRepository struct {
	db *sql.DB
}

func NewPostgresScheduledWorkoutRepository(db *sql.DB) irepo.ScheduledWorkoutRepository {
	return &PostgresScheduledWorkoutRepository{db: db}
}

func (r *PostgresScheduledWorkoutRepository) Create(ctx context.Context, sw *domain.ScheduledWorkout) error {
	if sw == nil {
		return fmt.Errorf("create scheduled workout: scheduled workout is nil")
	}

	const q = `
		INSERT INTO scheduled_workouts (user_id, workout_plan_id, scheduled_date)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, workout_plan_id, scheduled_date) DO NOTHING
		RETURNING id, created_at
	`

	if err := r.db.QueryRowContext(ctx, q, sw.UserID, sw.WorkoutPlanID, sw.ScheduledDate).Scan(&sw.ID, &sw.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return fmt.Errorf("create scheduled workout: %w", err)
	}

	return nil
}

func (r *PostgresScheduledWorkoutRepository) GetByUserAndDate(ctx context.Context, userID string, date time.Time) ([]domain.ScheduledWorkout, error) {
	const q = `
		SELECT id, user_id, workout_plan_id, scheduled_date, created_at
		FROM scheduled_workouts
		WHERE user_id = $1 AND scheduled_date = $2
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, q, userID, date)
	if err != nil {
		return nil, fmt.Errorf("get schedules by user and date: %w", err)
	}
	defer rows.Close()

	out := make([]domain.ScheduledWorkout, 0)
	for rows.Next() {
		var sw domain.ScheduledWorkout
		if err := rows.Scan(&sw.ID, &sw.UserID, &sw.WorkoutPlanID, &sw.ScheduledDate, &sw.CreatedAt); err != nil {
			return nil, fmt.Errorf("get schedules by user and date: %w", err)
		}
		out = append(out, sw)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get schedules by user and date: %w", err)
	}

	return out, nil
}

func (r *PostgresScheduledWorkoutRepository) GetByUser(ctx context.Context, userID string) ([]domain.ScheduledWorkout, error) {
	const q = `
		SELECT id, user_id, workout_plan_id, scheduled_date, created_at
		FROM scheduled_workouts
		WHERE user_id = $1
		ORDER BY scheduled_date ASC, created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("get schedules by user: %w", err)
	}
	defer rows.Close()

	out := make([]domain.ScheduledWorkout, 0)
	for rows.Next() {
		var sw domain.ScheduledWorkout
		if err := rows.Scan(&sw.ID, &sw.UserID, &sw.WorkoutPlanID, &sw.ScheduledDate, &sw.CreatedAt); err != nil {
			return nil, fmt.Errorf("get schedules by user: %w", err)
		}
		out = append(out, sw)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get schedules by user: %w", err)
	}

	return out, nil
}

func (r *PostgresScheduledWorkoutRepository) Delete(ctx context.Context, id string, userID string) error {
	const q = `
		DELETE FROM scheduled_workouts
		WHERE id = $1 AND user_id = $2
	`

	res, err := r.db.ExecContext(ctx, q, id, userID)
	if err != nil {
		return fmt.Errorf("delete schedule: %w", err)
	}

	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
