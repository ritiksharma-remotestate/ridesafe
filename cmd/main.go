package main

import (
	"net/http"
	"os"
	"os/signal"
	"ridesafe/database"
	"ridesafe/server"
	"ridesafe/utils"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const shutDownTimeOut = 10 * time.Second

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No .env file found")
	}
	utils.InitJWTSecret()
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if err := database.ConnectAndMigrate(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SSLMODE")); err != nil {
		logrus.Panicf("Failed to initialize and migrate databse and error: %+v", err)
	}
	logrus.Print("migration successful!!")

	srv := server.SetupRoutes()

	go func() {
		if err := srv.Run(":8080"); err != nil && err != http.ErrServerClosed {
			logrus.Panicf("failed to run server with error: %+v", err)
		}
	}()

	logrus.Print("server started at : 8080")
	<-done
	logrus.Info("shutting down server")
	if err := database.ShutdownDatabase(); err != nil {
		logrus.WithError(err).Error("Failed to close the database Connection")
	}
	if err := srv.Shutdown(shutDownTimeOut); err != nil {
		logrus.WithError(err).Panicf("FAILED TO GRACEFULLY SHUTDOWN SERVER")

	}
}
