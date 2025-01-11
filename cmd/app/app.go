package app

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"vchan.in/heyecho/handlers"
)

func Server() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.HideBanner = true // Hide the echo server banner to avoid server version disclosure in logs

	// Root level middleware
	if os.Getenv("LOG_LEVEL") == "debug" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Secure())                                                         // Use secure middleware to set security headers
	e.Use(middleware.Recover())                                                        // Recover middleware recovers from panics anywhere in the chain
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(5)))) // Rate limiter middleware
	e.Use(handlers.HandleAuthMiddleware)                                               // Auth middleware

	// Test page
	e.GET("/foo", handlers.HandleFoo)

	// Authencation pages with custom rate limiter 1 request per second only
	e.GET("/login", handlers.HandleLogin, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(1))))
	e.POST("/login", handlers.HandleLogin, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(1))))
	e.GET("/logout", handlers.HandleLogout)

	// Dashboard page
	e.GET("/", handlers.HandleDashboard)

	// Admin pages
	e.GET("/admin", handlers.HandleAdmin)
	e.GET("/admin/users", handlers.HandleUsers)
	e.POST("/admin/users", handlers.HandleCreateUser)
	e.PUT("/admin/users", handlers.HandleUpdateUser)

	httpAddr := os.Getenv("HTTP_ADDRESS")
	slog.Info("Server started at 0.0.0.0" + httpAddr)
	server := e.Start(httpAddr)
	if server != nil {
		slog.Error(server.Error())
	}
}
