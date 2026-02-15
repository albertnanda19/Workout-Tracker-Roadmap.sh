package migration

import (
	"database/sql"
	"errors"
	"os"
)

func RunMigrations(db *sql.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	var exists bool
	if err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'users'
		)
	`).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		b, err := os.ReadFile("schema.sql")
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(b)); err != nil {
			return err
		}
	}

	if _, err := db.Exec(`
		ALTER TABLE IF EXISTS scheduled_workouts
			ADD COLUMN IF NOT EXISTS user_id UUID,
			ADD COLUMN IF NOT EXISTS scheduled_date DATE;

		ALTER TABLE IF EXISTS scheduled_workouts
			ALTER COLUMN id SET DEFAULT gen_random_uuid();

		UPDATE scheduled_workouts sw
		SET user_id = wp.user_id
		FROM workout_plans wp
		WHERE sw.workout_plan_id = wp.id AND sw.user_id IS NULL;

		UPDATE scheduled_workouts
		SET scheduled_date = scheduled_at::date
		WHERE scheduled_date IS NULL AND scheduled_at IS NOT NULL;

		ALTER TABLE scheduled_workouts
			ALTER COLUMN user_id SET NOT NULL,
			ALTER COLUMN scheduled_date SET NOT NULL;

		ALTER TABLE scheduled_workouts
			ADD CONSTRAINT IF NOT EXISTS scheduled_workouts_user_id_fkey
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

		CREATE UNIQUE INDEX IF NOT EXISTS scheduled_workouts_unique
		ON scheduled_workouts(user_id, workout_plan_id, scheduled_date);

		CREATE INDEX IF NOT EXISTS idx_scheduled_user_date
		ON scheduled_workouts(user_id, scheduled_date);
	`); err != nil {
		return err
	}

	return nil
}
