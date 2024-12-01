package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()

	e.GET("/foo", handlers.HandleFoo)

	httpAddr := os.Getenv("HTTP_ADDRESS")
	slog.Info("Server started at " + httpAddr)
	e.Logger.Fatal(e.Start(httpAddr))
}
