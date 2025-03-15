package handlers

import (
	"encoding/base64"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/db"
	"vchan.in/heyecho/utils"
	"vchan.in/heyecho/views/auth/login"
)

// HandleLogin handles the login page.
func HandleLogin(c echo.Context) error {
	if c.Request().Method == "POST" {
		// Handle the login form submission
		email := c.FormValue("email")
		password := c.FormValue("password")

		slog.Info("Login attempt", "email", email)

		// Get user by email only
		dbuser, err := db.CallQueries.GetUserByEmail(c.Request().Context(), email)
		if err != nil {
			// If the user does not exist, send an error message to frontend
			color := "error"
			alert := "Incorrect email or password"
			return c.Redirect(http.StatusFound, "/login?level="+color+"&alert="+alert)
		}

		// Verify the password against the stored hash
		if utils.VerifyPasswordHash(password, dbuser.Password) {
			// If authentication is successful, set the cookie
			token := basicAuth(email, dbuser.Password)
			c.Response().Header().Set("Set-Cookie", "token="+token+"; Path=/")
			slog.Info("Login successful", "for email", email)

			// Set the user role for the session
			role = dbuser.Role
			c.Set("role", role)

			return c.Redirect(http.StatusFound, "/")
		} else {
			// Password doesn't match
			slog.Info("Login failed - incorrect password", "for email", email) 
			color := "error"
			alert := "Incorrect email or password"
			return c.Redirect(http.StatusFound, "/login?level="+color+"&alert="+alert)
		}
	}
	return render(c, login.Login(role))
}

// HandleLogout handles the logout page.
func HandleLogout(c echo.Context) error {
	c.Response().Header().Set("Set-Cookie", "token=; Path=/; Max-Age=-1")
	role = "guest"
	c.Set("role", role)
	return c.Redirect(http.StatusFound, "/login")
}

// basicAuth generates a basic auth token
func basicAuth(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}
