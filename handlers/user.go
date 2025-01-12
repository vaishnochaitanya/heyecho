package handlers

import (
	"crypto/subtle"
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

		// Check if the user email exists
		dbuser, err := db.CallQueries.GetUserByEmail(c.Request().Context(), email)
		if err != nil {
			// If the user does not exist, send an error message to frontend
			color := "error"
			alert := "Incorrect email or password"
			return c.Redirect(http.StatusFound, "/login?level="+color+"&alert="+alert)
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			slog.Error(err.Error())
			return c.Redirect(http.StatusFound, "/login")
		}

		// Check if the user email has a valid password
		if dbuser.Email != email && dbuser.Password != hashedPassword {
			slog.Info(email + " failed to login")
			return c.Redirect(http.StatusFound, "/login")
		}

		// Check if the provided credentials are valid
		if subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(dbuser.Password)) == 1 {
			// If authentication is successful, call the next handler
			c.Response().Header().Set("Set-Cookie", "token="+basicAuth(email, hashedPassword)+"; Path=/")
			slog.Info(email + " logged in")
			return c.Redirect(http.StatusFound, "/")
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
