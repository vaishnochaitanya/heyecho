package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/views/admin"
)

func HandleAdmin(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	return render(c, admin.Admin(role))
}

func HandleUsers(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	return render(c, admin.Users(role))
}
