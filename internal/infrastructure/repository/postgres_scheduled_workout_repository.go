package repository

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *PostgresScheduledWorkoutRepository) GetByUser(ctx context.Context, userID string, pagination domain.Pagination, filters domain.ScheduledWorkoutFilter) (domain.PaginatedResult[domain.ScheduledWorkout], error) {
	offset := (pagination.Page - 1) * pagination.Limit

	var date interface{} = nil
	if filters.Date != nil {
		date = *filters.Date
	}

	const countQ = `
		SELECT COUNT(1)
		FROM scheduled_workouts
		WHERE user_id = $1
		AND ($2::date IS NULL OR scheduled_date = $2)
	`

	var total int
	if err := r.db.QueryRowContext(ctx, countQ, userID, date).Scan(&total); err != nil {
		return domain.PaginatedResult[domain.ScheduledWorkout]{}, fmt.Errorf("get schedules by user: %w", err)
	}

	const q = `
		SELECT id, user_id, workout_plan_id, scheduled_date, created_at
		FROM scheduled_workouts
		WHERE user_id = $1
		AND ($2::date IS NULL OR scheduled_date = $2)
		ORDER BY scheduled_date DESC, created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, q, userID, date, pagination.Limit, offset)
	if err != nil {
		return domain.PaginatedResult[domain.ScheduledWorkout]{}, fmt.Errorf("get schedules by user: %w", err)
	}
	defer rows.Close()

	out := make([]domain.ScheduledWorkout, 0)
	for rows.Next() {
		var sw domain.ScheduledWorkout
		if err := rows.Scan(&sw.ID, &sw.UserID, &sw.WorkoutPlanID, &sw.ScheduledDate, &sw.CreatedAt); err != nil {
			return domain.PaginatedResult[domain.ScheduledWorkout]{}, fmt.Errorf("get schedules by user: %w", err)
		}
		out = append(out, sw)
	}
	if err := rows.Err(); err != nil {
		return domain.PaginatedResult[domain.ScheduledWorkout]{}, fmt.Errorf("get schedules by user: %w", err)
	}

	return domain.NewPaginatedResult(out, total, pagination), nil
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
