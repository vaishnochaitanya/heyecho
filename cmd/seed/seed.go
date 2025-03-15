package seed

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"vchan.in/heyecho/db"
	"vchan.in/heyecho/utils"
)

func run() error {
	ctx := context.Background()
	// create an admin user
	var UserUuid = strings.ReplaceAll(uuid.New().String(), "-", "")
	var UserName = "Admin"
	var UserEmail = "admin@example.com"
	var generatedPassphrase = utils.GeneratePassphrase
	var UserPassword, _ = utils.HashPassword(generatedPassphrase())
	var UserRole = "admin"

	err := db.CallQueries.CreateUser(ctx, db.CreateUserParams{
		Uuid:      UserUuid,
		Name:      UserName,
		Email:     UserEmail,
		Password:  UserPassword,
		Role:      UserRole,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		slog.Warn("Admin user creation failed: Already exists")
	} else {
		slog.Info("Admin user created with email: " + UserEmail + " and generated password: " + generatedPassphrase())
	}

	return nil
}

// Seed the database
func Seed() {
	if err := run(); err != nil {
		slog.Error(err.Error())
	}
}
