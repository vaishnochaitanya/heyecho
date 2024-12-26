package handlers

import (
	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/views/dashboard"
)

func HandleDashboard(c echo.Context) error {
	profile := c.Get("profile").(string)
	if profile == "" {
		profile = "guest"
	}
	return render(c, dashboard.Dashboard(profile))
}
