package repository

import (
	"context"
	"database/sql"
	"fmt"

	irepo "workout-tracker/internal/repository"
)

type PostgresWorkoutPlanChecker struct {
	db *sql.DB
}

func NewPostgresWorkoutPlanChecker(db *sql.DB) irepo.WorkoutPlanChecker {
	return &PostgresWorkoutPlanChecker{db: db}
}

func (c *PostgresWorkoutPlanChecker) GetOwnerID(ctx context.Context, workoutPlanID string) (string, error) {
	const q = `
		SELECT user_id
		FROM workout_plans
		WHERE id = $1
	`

	var userID string
	if err := c.db.QueryRowContext(ctx, q, workoutPlanID).Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", fmt.Errorf("get workout plan owner: %w", err)
	}

	return userID, nil
}
