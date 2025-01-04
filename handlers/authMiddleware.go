package handlers

import (
	"crypto/subtle"
	"encoding/base64"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/db"
	"vchan.in/heyecho/types"
)

var role = types.Role

// Basic authentication middleware
func HandleAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Add exception for the following paths
		excepted, nextPath := isException(c, next)
		if excepted {
			return nextPath // Continue to the next handler
		}

		// Extract the token from the cookie
		cookie, ok := c.Cookie("token")
		if ok != nil {
			return c.Redirect(http.StatusFound, "/login")
		}

		// Decode the token to get the email and password
		token, ok := base64.StdEncoding.DecodeString(cookie.Value)
		if ok != nil {
			return c.Redirect(http.StatusFound, "/login")
		}

		// Split the token into email and password
		credentials := string(token)
		email, password, _ := strings.Cut(credentials, ":")

		// Check if the user email exists
		user, err := db.CallQueries.GetUserByEmail(c.Request().Context(), email)
		if err != nil {
			// If the user does not exist, redirect to the login page
			return c.Redirect(http.StatusFound, "/login")
		}
		// Check if the user email has a valid password
		if user.Email != email && user.Password != password {
			slog.Info(email + " failed to login")
			return c.Redirect(http.StatusFound, "/login")
		}

		// Check if the provided credentials are valid
		if subtle.ConstantTimeCompare([]byte(password), []byte(user.Password)) == 1 {
			// If authentication is successful, set the role
			if email == "admin@example.com" {
				role = "admin"
			} else if email != "" {
				role = "user"
			}
			c.Set("role", role)
			return next(c)
		}

		// If authentication fails, set the role to guest
		return next(c)
	}
}

func isException(c echo.Context, next echo.HandlerFunc) (bool, error) {
	exceptionPaths := []string{
		"/login",
		"/logout",
	}
	for _, path := range exceptionPaths {
		if c.Path() == path {
			return true, next(c)
		}
	}
	return false, nil
}
