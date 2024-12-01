package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()

	e.GET("/foo", handleFoo)

	httpAddr := os.Getenv("HTTP_ADDRESS")
	slog.Info("Server started at " + httpAddr)
	e.Logger.Fatal(e.Start(httpAddr))
}

func handleFoo(c echo.Context) error {
	return c.String(http.StatusOK, "foo")
}
