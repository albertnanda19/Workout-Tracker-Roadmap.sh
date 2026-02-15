package seeder

import (
	"database/sql"
	"errors"
	"time"
)

type exerciseSeed struct {
	Name        string
	Category    string
	MuscleGroup string
}

func RunSeeders(db *sql.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	seeds := []exerciseSeed{
		{Name: "Bench Press", Category: "strength", MuscleGroup: "chest"},
		{Name: "Squat", Category: "strength", MuscleGroup: "legs"},
		{Name: "Deadlift", Category: "strength", MuscleGroup: "back"},
		{Name: "Pull Up", Category: "strength", MuscleGroup: "back"},
		{Name: "Push Up", Category: "strength", MuscleGroup: "chest"},
		{Name: "Lunges", Category: "strength", MuscleGroup: "legs"},
		{Name: "Plank", Category: "strength", MuscleGroup: "core"},
		{Name: "Shoulder Press", Category: "strength", MuscleGroup: "shoulders"},
		{Name: "Bicep Curl", Category: "strength", MuscleGroup: "arms"},
		{Name: "Tricep Dip", Category: "strength", MuscleGroup: "arms"},
		{Name: "Running", Category: "cardio", MuscleGroup: "legs"},
		{Name: "Cycling", Category: "cardio", MuscleGroup: "legs"},
		{Name: "Jump Rope", Category: "cardio", MuscleGroup: "core"},
		{Name: "Leg Press", Category: "strength", MuscleGroup: "legs"},
		{Name: "Lat Pulldown", Category: "strength", MuscleGroup: "back"},
		{Name: "Chest Fly", Category: "strength", MuscleGroup: "chest"},
		{Name: "Leg Curl", Category: "strength", MuscleGroup: "legs"},
		{Name: "Leg Extension", Category: "strength", MuscleGroup: "legs"},
		{Name: "Russian Twist", Category: "flexibility", MuscleGroup: "core"},
		{Name: "Mountain Climbers", Category: "cardio", MuscleGroup: "core"},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(`
		INSERT INTO exercises (name, category, muscle_group)
		VALUES ($1, $2, $3)
		ON CONFLICT (name) DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range seeds {
		if _, err := stmt.Exec(s.Name, s.Category, s.MuscleGroup); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if err := seedScheduledWorkout(db); err != nil {
		return err
	}

	return nil
}

func seedScheduledWorkout(db *sql.DB) error {
	var userID string
	if err := db.QueryRow(`
		SELECT id
		FROM users
		ORDER BY created_at ASC
		LIMIT 1
	`).Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	var planID string
	if err := db.QueryRow(`
		SELECT id
		FROM workout_plans
		WHERE user_id = $1
		ORDER BY created_at ASC
		LIMIT 1
	`, userID).Scan(&planID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	date := time.Now().UTC().Add(24 * time.Hour)
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	if _, err := db.Exec(`
		INSERT INTO scheduled_workouts (user_id, workout_plan_id, scheduled_date)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, workout_plan_id, scheduled_date) DO NOTHING
	`, userID, planID, date); err != nil {
		return err
	}

	return nil
}
