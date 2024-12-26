package app

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/handlers"
)

func Server() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.HideBanner = true // Hide the echo server banner to avoid server version disclosure in logs

	e.GET("/foo", handlers.HandleFoo)

	// Authencation pages
	e.GET("/login", handlers.HandleLogin)

	// Dashboard page
	e.GET("/", handlers.HandleDashboard)

	httpAddr := os.Getenv("HTTP_ADDRESS")
	slog.Info("Server started at " + httpAddr)
	e.Logger.Fatal(e.Start(httpAddr))
}
