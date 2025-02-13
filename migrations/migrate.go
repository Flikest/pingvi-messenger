package migrations

import (
	"database/sql"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func CreateMigrations(db *sql.DB, filePath string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Debug("error creating driver")
	}

	migration, err := migrate.NewWithDatabaseInstance(
		filePath,
		"postgres", driver)
	if err != nil {
		slog.Debug("error during migration")
	}

	migration.Up()
	slog.Debug("migrations completed")
}
