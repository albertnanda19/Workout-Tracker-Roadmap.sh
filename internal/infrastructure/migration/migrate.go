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

	if exists {
		return nil
	}

	b, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(b)); err != nil {
		return err
	}

	return nil
}
