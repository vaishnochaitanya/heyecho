package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"vchan.in/heyecho/db"
	"vchan.in/heyecho/logger"
	"vchan.in/heyecho/utils"
)

func CreateNewAdminUser(email string) (bool, error) {
	// Check if the email is empty
	if email == "" {
		return false, fmt.Errorf("email is empty")
	}

	// Check if the user already exists
	ctx := context.Background()
	_, err := db.CallQueries.GetUserByEmail(ctx, email)
	if err == nil {
		return false, fmt.Errorf("user already exists")
	}

	// Custom logger
	level := slog.LevelInfo // default level
	slogger := slog.New(logger.NewSimpleHandler(os.Stdout, level))
	slog.SetDefault(slogger)

	// create an admin user
	var UserUuid = strings.ReplaceAll(uuid.New().String(), "-", "")
	var UserName = email
	var UserEmail = email
	var generatedPassphrase = utils.GeneratePassphrase()
	var UserPassword, _ = utils.HashPassword(generatedPassphrase)
	var UserRole = "admin"

	err = db.CallQueries.CreateUser(ctx, db.CreateUserParams{
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
		slog.Info("Admin user created with email: " + UserEmail + " and generated password: " + generatedPassphrase)
	}

	return true, nil
}

func ResetAdminPassword(email string) (bool, error) {
	// Check if the email is empty
	if email == "" {
		return false, fmt.Errorf("email is empty")
	}

	// Check if the user exists
	ctx := context.Background()
	dbuser, err := db.CallQueries.GetUserByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	// Generate a new password
	newPassword := uuid.New().String()
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return false, fmt.Errorf("error hashing password: %w", err)
	}

	// Update the password in the database
	_, err = db.CallQueries.UpdateUser(ctx, db.UpdateUserParams{
		Uuid:      dbuser.Uuid,
		Password:  hashedPassword,
		UpdatedAt: dbuser.UpdatedAt,
	})
	if err != nil {
		return false, fmt.Errorf("error updating password: %w", err)
	} else {
		slog.Info("Admin Argument Reset Password: Password reset", "Email:", email, "New Password:", newPassword)
	}
	return true, nil
}
