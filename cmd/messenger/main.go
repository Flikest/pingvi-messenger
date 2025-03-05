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
	"github.com/Flikest/PingviMessenger/pkg/logger"
	"github.com/Flikest/PingviMessenger/rabbitmq"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	log := logger.InitLogger(os.Getenv("LVL_DEPLOYMENT"))

	db, err := postgresql.NewDatabase(&postgresql.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		log.Info("error connecting to database")
	}

	migrations.CreateMigrations(db, "file://pingviMessenger/migrations/sql/")

	storage := storage.NewStorage(db, context.Background())
	services := services.NewServices(storage)
	handler := handler.NewHandler(services)
	router := handler.InitRouter()

	if err := router.Run(":9000"); err != nil {
		log.Info("server not runing")
	}
	log.Info("server is starting")

	rabbitmq.Consume(db)
}
