package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"image/png"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/disintegration/letteravatar"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/db"
	"vchan.in/heyecho/types"
	"vchan.in/heyecho/views/admin"
)

func generatedPhoto(name string) string {
	firstLetter := rune(name[0]) // Get first letter of the name and generate a letter avatar
	img, err := letteravatar.Draw(40, firstLetter, nil)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return ""
	}
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func HandleAdmin(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	return render(c, admin.Admin(role))
}

// List all users	GET /admin/users
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
			Uuid:           u.Uuid,
			Name:           u.Name,
			Email:          u.Email,
			Role:           u.Role,
			CreatedAt:      u.CreatedAt.Format("7 Jan 2006 10:36 AM"),
			UpdatedAt:      u.UpdatedAt.Format("7 Jan 2006 10:36 AM"),
			GeneratedPhoto: generatedPhoto(u.Name),
		})
	}
	return render(c, admin.Users(role, users))
}

// Create a new user	POST /admin/users
func HandleCreateUser(c echo.Context) error {
	var user struct {
		Name     string `json:"name" form:"name"`
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
		Role     string `json:"role" form:"role"`
	}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Only if the current user role is admin
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	slog.Info("Creating user: ", user)

	if user.Name == "" || user.Email == "" || user.Password == "" || user.Role == "" {
		return c.JSON(http.StatusBadRequest, "Name, Email, Password and Role are required")
	}

	ctx := context.Background()
	err := db.CallQueries.CreateUser(ctx, db.CreateUserParams{
		Uuid:      strings.ReplaceAll(uuid.New().String(), "-", ""),
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	timeNow := time.Now().Format("7 Jan 2006 10:36 AM")
	GeneratedPhoto := generatedPhoto(user.Name)
	responseHTML := `<tr
                        class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
                        <th scope="row" class="flex items-center px-6 py-4 text-gray-900 whitespace-nowrap dark:text-white">
                            <img class="w-10 h-10 rounded-full" src=` + GeneratedPhoto + ` alt="avatar">
                            <div class="ps-3">
                                <div id="name" class="text-base font-semibold">` + user.Name + `</div>
                                <div id="email" class="font-normal text-gray-500">` + user.Email + `</div>
                            </div>
                        </th>
                        <td id="role" class="px-6 py-4">
                            ` + user.Role + `
                        </td>
                        <td class="px-6 py-4">
                            <div class="flex items-center">
                                ` + timeNow + `
                            </div>
                        </td>
                        <td class="px-6 py-4">
                            <!-- Modal toggle -->
                            <a href="#" id="editUser" type="button"
                                data-modal-target="userModal" data-modal-show="userModal"
                                class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit user</a>
                        </td>
                    </tr>`
	return c.HTML(http.StatusOK, responseHTML)
}

// Update a user	PUT /admin/users
func HandleUpdateUser(c echo.Context) error {
	var user struct {
		Uuid  string `json:"uuid" form:"uuid"`
		Name  string `json:"name" form:"name"`
		Email string `json:"email" form:"email"`
		Role  string `json:"role" form:"role"`
	}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Only if the current user role is admin
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, types.Response{Status: "error", Message: "Forbidden"})
	}

	ctx := context.Background()
	// Check if the user exists
	dbUser, err := db.CallQueries.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if dbUser.Uuid == "" {
		return c.JSON(http.StatusNotFound, types.Response{Status: "error", Message: "User not found"})
	}

	// Update the user
	_, err = db.CallQueries.UpdateUser(ctx, db.UpdateUserParams{
		Uuid:      user.Uuid,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Response{Status: "error", Message: err.Error()})
	}
	timeNow := time.Now().Format("7 Jan 2006 10:36 AM")
	responseHTML := `<tr
                        class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
                        <th scope="row" class="flex items-center px-6 py-4 text-gray-900 whitespace-nowrap dark:text-white">
                            <img class="w-10 h-10 rounded-full" src=` + generatedPhoto(user.Name) + ` alt="avatar">
                            <div class="ps-3">
                                <div id="name" class="text-base font-semibold">` + user.Name + `</div>
                                <div id="email" class="font-normal text-gray-500">` + user.Email + `</div>
                            </div>
                        </th>
                        <td id="role" class="px-6 py-4">
                            ` + user.Role + `
                        </td>
                        <td class="px-6 py-4">
                            <div class="flex items-center">
                                ` + timeNow + `
                            </div>
                        </td>
                        <td class="px-6 py-4">
                            <!-- Modal toggle -->
                            <a href="#" id="editUser" type="button"
                                data-modal-target="userModal" data-modal-show="userModal"
                                class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit user</a>
                        </td>
                    </tr>`
	return c.HTML(http.StatusOK, responseHTML)
}
