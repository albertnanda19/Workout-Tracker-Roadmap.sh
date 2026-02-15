package repository

import (
	"context"
	"database/sql"
	"fmt"

	"workout-tracker/internal/domain"
)

type PostgresExerciseRepository struct {
	db *sql.DB
}

func NewPostgresExerciseRepository(db *sql.DB) domain.ExerciseRepository {
	return &PostgresExerciseRepository{db: db}
}

func (r *PostgresExerciseRepository) GetAll(ctx context.Context) ([]domain.Exercise, error) {
	const q = `
		SELECT id, name, description, category, muscle_group
		FROM exercises
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("get all exercises: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Exercise, 0)
	for rows.Next() {
		var e domain.Exercise
		var description sql.NullString
		var category sql.NullString
		var muscleGroup sql.NullString
		if err := rows.Scan(&e.ID, &e.Name, &description, &category, &muscleGroup); err != nil {
			return nil, fmt.Errorf("get all exercises: %w", err)
		}
		e.Description = description.String
		e.Category = category.String
		e.MuscleGroup = muscleGroup.String
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get all exercises: %w", err)
	}

	return out, nil
}

func (r *PostgresExerciseRepository) GetByID(ctx context.Context, id string) (*domain.Exercise, error) {
	const q = `
		SELECT id, name, description, category, muscle_group
		FROM exercises
		WHERE id = $1
	`

	var e domain.Exercise
	var description sql.NullString
	var category sql.NullString
	var muscleGroup sql.NullString
	if err := r.db.QueryRowContext(ctx, q, id).Scan(&e.ID, &e.Name, &description, &category, &muscleGroup); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("get exercise by id: %w", err)
	}
	e.Description = description.String
	e.Category = category.String
	e.MuscleGroup = muscleGroup.String

	return &e, nil
}
