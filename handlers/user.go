package handlers

import (
	"crypto/subtle"
	"encoding/base64"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/views/auth/login"
)

// HandleLogin handles the login page.
func HandleLogin(c echo.Context) error {
	if c.Request().Method == "POST" {
		// Handle the login form submission
		email := c.FormValue("email")
		password := c.FormValue("password")

		// List of valid user credentials
		validCredentials := map[string]string{
			"admin@example.com": "admin",
			"user@example.com":  "user",
		}

		// Check if the provided credentials are valid
		if pass, ok := validCredentials[email]; ok && subtle.ConstantTimeCompare([]byte(password), []byte(pass)) == 1 {
			// If authentication is successful, call the next handler
			c.Response().Header().Set("Set-Cookie", "token="+basicAuth(email, password)+"; Path=/")
			slog.Info(email + " logged in")
			return c.Redirect(http.StatusFound, "/")
		}

		// If the credentials are incorrect, set the profile to guest
		slog.Info(email + " failed to login")
		return c.Redirect(http.StatusFound, "/login")
	}
	return render(c, login.Login(profile))
}

// HandleLogout handles the logout page.
func HandleLogout(c echo.Context) error {
	c.Response().Header().Set("Set-Cookie", "token=; Path=/; Max-Age=-1")
	profile = "guest"
	c.Set("profile", profile)
	return c.Redirect(http.StatusFound, "/login")
}

// basicAuth generates a basic auth token
func basicAuth(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}
