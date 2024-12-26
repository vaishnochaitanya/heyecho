package handlers

import (
	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/views/auth/login"
)

func HandleLogin(c echo.Context) error {
	return render(c, login.Login(profile))
}
