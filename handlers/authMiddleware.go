package handlers

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/types"
)

var profile = types.Profile

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

		// List of valid user credentials
		validCredentials := map[string]string{
			"admin@example.com": "admin",
			"user@example.com":  "user",
		}

		// Check if the provided credentials are valid
		if pass, ok := validCredentials[email]; ok && subtle.ConstantTimeCompare([]byte(password), []byte(pass)) == 1 {
			// If authentication is successful, set the profile
			if email == "admin@example.com" {
				profile = "admin"
			} else if email != "" {
				profile = "user"
			}
			c.Set("profile", profile)
			return next(c)
		}

		// If authentication fails, set the profile to guest
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
