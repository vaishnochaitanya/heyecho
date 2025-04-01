package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/db"
	"vchan.in/heyecho/utils"
)

// CLI argument handler to reset the admin password
func HandleAdminArgumentResetPassword(c echo.Context) error {
	// Get the email from the CLI argument
	email := c.Param("email")
	if email == "" {
		slog.Error("Email is empty")
		return c.JSON(http.StatusBadRequest, "Email is empty")
	}

	// Check if the user exists
	ctx := context.Background()
	dbuser, err := db.CallQueries.GetUserByEmail(ctx, email)
	if err != nil {
		slog.Error("User not found", "error", err.Error())
		return c.JSON(http.StatusNotFound, "User not found")
	}

	// Generate a new password
	newPassword := uuid.New().String()
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		slog.Error("Admin Argument Reset Password: Error hashing password", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, "Admin Argument Reset Password: Error hashing password")
	}

	// Update the password in the database
	_, err = db.CallQueries.UpdateUser(ctx, db.UpdateUserParams{
		Uuid:      dbuser.Uuid,
		Password:  hashedPassword,
		UpdatedAt: dbuser.UpdatedAt,
	})
	if err != nil {
		slog.Error("Admin Argument Reset Password: Error updating password", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, "Admin Argument Reset Password: Error updating password")
	}

	slog.Info("Admin Argument Reset Password: Password reset", "Email:", email, "New Password:", newPassword)
	return c.JSON(http.StatusOK, newPassword)
}

// CLI argument handler to create a new admin user. The email is passed as an argument in the CLI
func HandleAdminArgumentCreateUser(c echo.Context) error {
	// Get the email from the CLI argument
	requestEmail := c.Param("email")
	if requestEmail == "" {
		slog.Error("Email is empty")
		return c.JSON(http.StatusBadRequest, "Email is empty")
	}

	// Check if the user exists
	ctx := context.Background()
	_, err := db.CallQueries.GetUserByEmail(ctx, requestEmail)
	if err == nil {
		slog.Error("User already exists", "email", requestEmail)
		return c.JSON(http.StatusConflict, "User already exists")
	}

	// create an admin user
	var UserUuid = strings.ReplaceAll(uuid.New().String(), "-", "")
	var UserName = "Admin"
	var UserEmail = requestEmail
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
		return c.JSON(http.StatusConflict, "Admin user creation failed: Already exists")
	}

	slog.Info("Admin user created with email: " + UserEmail + " and generated password: " + generatedPassphrase)
	return c.JSON(http.StatusOK, "Admin user created with email: "+UserEmail+" and generated password: "+generatedPassphrase)
}
