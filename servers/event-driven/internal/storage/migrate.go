package storage

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB, migrationsDir string) error {
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("RunMigrations: failed to make driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "pgx", driver)
	if err != nil {
		return fmt.Errorf("RunMigrations: failed to make migrator: %w", err)
	}

	err = m.Up()
	if !(err == nil || err == migrate.ErrNoChange) {
		return fmt.Errorf("RunMigrations: failed to up migrations: %w", err)
	}

	return nil
}
