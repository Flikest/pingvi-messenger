package main

import (
	"context"
	"log/slog"
	"os"

	services "github.com/Flikest/PingviMessenger/internal/controller"
	"github.com/Flikest/PingviMessenger/internal/handler"
	"github.com/Flikest/PingviMessenger/internal/storage"
	"github.com/Flikest/PingviMessenger/migrations"
	postgresql "github.com/Flikest/PingviMessenger/pkg/clientdb/postgresql"
	"github.com/Flikest/PingviMessenger/rabbitmq"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

// @title           Pingue API
// @version         1.0
// @description     This is pengvi messenger.
// @host      localhost:9000
func main() {
	err := godotenv.Load("C:/Users/User/Desktop/pingviMessenger/.env")
	if err != nil {
		slog.Info("error reading .env: ", err)
	}

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

	migrations.CreateMigrations(db, "file://C:/Users/User/Desktop/pingviMessenger/migrations/sql/")

	rabbitmq.Consume(db)

	storage := storage.NewStorage(db, context.Background())
	services := services.NewServices(storage)
	handler := handler.NewHandler(services)
	h := handler.InitRouter()

	if err := h.Run(":9000"); err != nil {
		log.Debug("server not started", err)
	} else {
		log.Info("server started on port:", "9000")
	}

}
