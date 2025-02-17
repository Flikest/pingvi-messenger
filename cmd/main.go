package main

import (
	"net/http"
	"os"

	"github.com/Flikest/PingviMessenger/internal/handler"
	"github.com/Flikest/PingviMessenger/internal/services"
	"github.com/Flikest/PingviMessenger/internal/storage"
	"github.com/Flikest/PingviMessenger/migrations"
	postgresql "github.com/Flikest/PingviMessenger/pkg/clientdb/postgresql"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db, err := postgresql.NewDatabase(&postgresql.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		panic("error connecting to database")
	}

	migrations.CreateMigrations(db, "file://migrations/sql/")

	s := storage.InitStorage(db)
	serv := services.InitService(s)
	Handler := handler.InitHandler(serv)
	h := Handler.InitRouter()

	if err := http.ListenAndServe(":9000", h); err != nil {
		log.Debug("server not started", err)
	} else {
		log.Info("server started on port:", "9000")
	}

}
