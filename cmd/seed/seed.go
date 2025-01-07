package seed

import (
	"context"
	"log/slog"
	"time"

	"vchan.in/heyecho/db"
)

func run() error {
	ctx := context.Background()
	// create an admin user
	err := db.CallQueries.CreateUser(ctx, db.CreateUserParams{
		Name:      "Admin",
		Email:     "admin@example.com",
		Password:  "admin",
		Role:      "admin",
		CreatedAt: time.Now(),
	})
	if err != nil {
		slog.Warn("Admin user creation failed: Already exists")
	} else {
		slog.Info("Admin user created with email: admin@example.com and password: admin")
	}

	return nil
}

// Seed the database
func Seed() {
	if err := run(); err != nil {
		slog.Error(err.Error())
	}
}
