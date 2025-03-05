package migrations

import (
	"database/sql"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func CreateMigrations(db *sql.DB, filePath string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Info("error during migration: ", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		filePath,
		"postgres", driver)
	if err != nil {
		slog.Info("error during migration: ", err)
	}

	err = migration.Up()
	if err != nil {
		slog.Info("migration dont up: ", err)
	} else {
		slog.Info("migrations completed")
	}
}
