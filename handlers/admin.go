package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/views/admin"
)

func HandleAdmin(c echo.Context) error {
	profile := c.Get("profile").(string)
	if profile != "admin" {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	return render(c, admin.Admin(profile))
}

func HandleUsers(c echo.Context) error {
	profile := c.Get("profile").(string)
	if profile != "admin" {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	return render(c, admin.Users(profile))
}
