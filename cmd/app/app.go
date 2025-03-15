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
	"vchan.in/heyecho/logger"
)

func Server() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file. Please create a .env file in the root directory")
	}

	level := slog.LevelInfo // default level

	// Set log level based on environment variable
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	// Create a custom simple logger with the desired format
	slogger := slog.New(logger.NewSimpleHandler(os.Stdout, level))
	slog.SetDefault(slogger)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Root level middleware
	if os.Getenv("LOG_LEVEL") == "debug" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.CORS())                                                           // Use CORS middleware to allow cross-origin requests
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))                  // Use Gzip middleware to compress the response
	e.Use(middleware.Secure())                                                         // Use secure middleware to set security headers
	e.Use(middleware.Recover())                                                        // Recover middleware recovers from panics anywhere in the chain
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(5)))) // Rate limiter middleware
	e.Use(handlers.HandleAuthMiddleware)                                               // Auth middleware

	// Test page
	e.GET("/foo", handlers.HandleFoo)

	// Authencation pages with custom rate limiter 5 request per second only
	e.GET("/login", handlers.HandleLogin, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(5))))
	e.POST("/login", handlers.HandleLogin, middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(1))))
	e.GET("/logout", handlers.HandleLogout)

	// Dashboard page
	e.GET("/", handlers.HandleDashboard)

	// Admin pages
	const adminUsersPath = "/admin/users"
	e.GET("/admin", handlers.HandleAdmin)
	e.GET(adminUsersPath, handlers.HandleUsers)
	e.POST(adminUsersPath, handlers.HandleCreateUser)
	e.PUT(adminUsersPath, handlers.HandleUpdateUser)

	httpAddr := os.Getenv("HTTP_ADDRESS")
	slog.Info("Server started at 0.0.0.0" + httpAddr)
	if e.Start(httpAddr) != nil {
		slog.Error("Server failed to start")
	}
}
