package handlers

import (
	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/views/dashboard"
)

func HandleDashboard(c echo.Context) error {
	return render(c, dashboard.Dashboard(profile))
}
