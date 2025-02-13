package main

import (
	"os"

	"github.com/Flikest/PingviMessenger/migrations"
	postgresql "github.com/Flikest/PingviMessenger/pkg/clientdb/postgresql"
)

func main() {

	db, err := postgresql.NewDatabase(&postgresql.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_DBNAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		panic("error connecting to database")
	}

	migrations.CreateMigrations(db, "file://migrations/sql/0001.messenger.up.sql")
}
