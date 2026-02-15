package repository

import (
	"context"
	"database/sql"
	"fmt"

	"workout-tracker/internal/domain"
)

type PostgresWorkoutRepository struct {
	db *sql.DB
}

func NewPostgresWorkoutRepository(db *sql.DB) domain.WorkoutRepository {
	return &PostgresWorkoutRepository{db: db}
}

func (r *PostgresWorkoutRepository) CreatePlan(ctx context.Context, plan *domain.WorkoutPlan, exercises []domain.WorkoutPlanExercise) error {
	if plan == nil {
		return fmt.Errorf("create plan: plan is nil")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create plan: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const insertPlan = `
		INSERT INTO workout_plans (user_id, name, notes)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var planID string
	if err := tx.QueryRowContext(ctx, insertPlan, plan.UserID, plan.Name, plan.Notes).Scan(&planID); err != nil {
		return fmt.Errorf("create plan: %w", err)
	}

	const insertPlanExercise = `
		INSERT INTO workout_plan_exercises (workout_plan_id, exercise_id, sets, reps, weight, order_index)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, ex := range exercises {
		if _, err := tx.ExecContext(ctx, insertPlanExercise, planID, ex.ExerciseID, ex.Sets, ex.Reps, ex.Weight, ex.OrderIndex); err != nil {
			return fmt.Errorf("create plan: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create plan: %w", err)
	}

	plan.ID = planID
	return nil
}

func (r *PostgresWorkoutRepository) UpdatePlan(ctx context.Context, plan *domain.WorkoutPlan, exercises []domain.WorkoutPlanExercise) error {
	if plan == nil {
		return fmt.Errorf("update plan: plan is nil")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("update plan: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var existingID string
	if err := tx.QueryRowContext(ctx, `
		SELECT id
		FROM workout_plans
		WHERE id = $1 AND user_id = $2
	`, plan.ID, plan.UserID).Scan(&existingID); err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return fmt.Errorf("update plan: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE workout_plans
		SET name = $1, notes = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4
	`, plan.Name, plan.Notes, plan.ID, plan.UserID); err != nil {
		return fmt.Errorf("update plan: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM workout_plan_exercises
		WHERE workout_plan_id = $1
	`, plan.ID); err != nil {
		return fmt.Errorf("update plan: %w", err)
	}

	const insertPlanExercise = `
		INSERT INTO workout_plan_exercises (workout_plan_id, exercise_id, sets, reps, weight, order_index)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, ex := range exercises {
		if _, err := tx.ExecContext(ctx, insertPlanExercise, plan.ID, ex.ExerciseID, ex.Sets, ex.Reps, ex.Weight, ex.OrderIndex); err != nil {
			return fmt.Errorf("update plan: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("update plan: %w", err)
	}

	return nil
}

func (r *PostgresWorkoutRepository) GetPlansByUser(ctx context.Context, userID string, pagination domain.Pagination, filters domain.WorkoutPlanFilter) (domain.PaginatedResult[domain.WorkoutPlan], error) {
	offset := (pagination.Page - 1) * pagination.Limit

	const countQ = `
		SELECT COUNT(1)
		FROM workout_plans
		WHERE user_id = $1
		AND ($2 = '' OR name ILIKE '%' || $2 || '%')
	`

	var total int
	if err := r.db.QueryRowContext(ctx, countQ, userID, filters.Name).Scan(&total); err != nil {
		return domain.PaginatedResult[domain.WorkoutPlan]{}, fmt.Errorf("get plans by user: %w", err)
	}

	const q = `
		SELECT id, user_id, name, notes, created_at, updated_at
		FROM workout_plans
		WHERE user_id = $1
		AND ($2 = '' OR name ILIKE '%' || $2 || '%')
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, q, userID, filters.Name, pagination.Limit, offset)
	if err != nil {
		return domain.PaginatedResult[domain.WorkoutPlan]{}, fmt.Errorf("get plans by user: %w", err)
	}
	defer rows.Close()

	out := make([]domain.WorkoutPlan, 0)
	for rows.Next() {
		var p domain.WorkoutPlan
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return domain.PaginatedResult[domain.WorkoutPlan]{}, fmt.Errorf("get plans by user: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return domain.PaginatedResult[domain.WorkoutPlan]{}, fmt.Errorf("get plans by user: %w", err)
	}

	return domain.NewPaginatedResult(out, total, pagination), nil
}

func (r *PostgresWorkoutRepository) GetPlanByID(ctx context.Context, id string, userID string) (*domain.WorkoutPlan, error) {
	const q = `
		SELECT id, user_id, name, notes, created_at, updated_at
		FROM workout_plans
		WHERE id = $1 AND user_id = $2
	`

	var p domain.WorkoutPlan
	if err := r.db.QueryRowContext(ctx, q, id, userID).Scan(&p.ID, &p.UserID, &p.Name, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get plan by id: %w", err)
	}

	return &p, nil
}

func (r *PostgresWorkoutRepository) DeletePlan(ctx context.Context, id string, userID string) error {
	const q = `
		DELETE FROM workout_plans
		WHERE id = $1 AND user_id = $2
	`

	res, err := r.db.ExecContext(ctx, q, id, userID)
	if err != nil {
		return fmt.Errorf("delete plan: %w", err)
	}

	_, _ = res.RowsAffected()
	return nil
}
