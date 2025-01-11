package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"image/png"
	"net/http"

	"github.com/disintegration/letteravatar"
	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/db"
	"vchan.in/heyecho/types"
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
	// List all users from the database ListUsers
	ctx := context.Background()
	dbUsers, err := db.CallQueries.ListUsers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	users := []types.User{}
	for _, u := range dbUsers {
		users = append(users, types.User{
			Uuid:      u.Uuid,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt.Format("7 Jan 2006 10:36 AM"),
			UpdatedAt: u.UpdatedAt.Format("7 Jan 2006 10:36 AM"),
			GeneratedPhoto: func() string {
				img, err := letteravatar.Draw(100, 'A', nil)
				if err != nil {
					return ""
				}
				var buf bytes.Buffer
				if err := png.Encode(&buf, img); err != nil {
					return ""
				}
				return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
			}(),
		})
	}
	return render(c, admin.Users(role, users))
}

// API
func HandleAPIUsers(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
	return c.JSON(http.StatusOK, "API Users")
}
