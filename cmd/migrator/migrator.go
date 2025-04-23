package main

import (
	"database/sql"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DB_PATH"))
	if err != nil {
		panic("not found db")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic("failed to create driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations/sql",
		"postgres", driver)
	m.Up()
}
