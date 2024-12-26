package handlers

import (
	"crypto/subtle"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Basic authentication middleware
func HandleAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	// Middleware function to handle authentication
	return func(c echo.Context) error {
		// Implement basic authentication here
		username, password, ok := c.Request().BasicAuth()
		if !ok {
			// If the request does not contain basic auth headers, redirect to the login page
			return c.Redirect(http.StatusFound, "/login")
		}

		// List of valid user credentials
		validCredentials := map[string]string{
			"admin": "admin",
			"user":  "user",
		}

		// Check if the provided credentials are valid
		if pass, ok := validCredentials[username]; ok && subtle.ConstantTimeCompare([]byte(password), []byte(pass)) == 1 {
			// If authentication is successful, call the next handler
			return next(c)
		}

		// If authentication fails, redirect to the login page
		return c.Redirect(http.StatusFound, "/login")
	}
}
