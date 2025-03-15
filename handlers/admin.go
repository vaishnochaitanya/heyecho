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
	"vchan.in/heyecho/utils"
	"vchan.in/heyecho/views/admin"
)

const dateFormat = "7 Jan 2006 10:36 AM"

func generatedPhoto(name string) string {
	firstLetter := rune(name[0]) // Get first letter of the name and generate a letter avatar
	img, err := letteravatar.Draw(40, firstLetter, nil)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		slog.Error("Error encoding image", "error", err.Error())
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
			CreatedAt:      u.CreatedAt.Format(dateFormat),
			UpdatedAt:      u.UpdatedAt.Format(dateFormat),
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

	slog.Info("Creating user: ", slog.Any("user", user))

	if user.Name == "" || user.Email == "" || user.Password == "" || user.Role == "" {
		return c.JSON(http.StatusBadRequest, "Name, Email, Password and Role are required")
	}

	ctx := context.Background()

	// Continue if the user does not exist
	_, err := db.CallQueries.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return c.JSON(http.StatusConflict, "User already exists")
	}

	// Create the user
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		slog.Debug("Add User: Error hashing password", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	uniqueId := strings.ReplaceAll(uuid.New().String(), "-", "")

	if err = db.CallQueries.CreateUser(ctx, db.CreateUserParams{
		Uuid:      uniqueId,
		Name:      user.Name,
		Email:     user.Email,
		Password:  hashedPassword,
		Role:      user.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		slog.Error("Error creating user", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	timeNow := time.Now().Format(dateFormat)
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
		Uuid     string `json:"uuid" form:"uuid"`
		Name     string `json:"name" form:"name"`
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
		Role     string `json:"role" form:"role"`
	}
	user.Password = "123456@temp!" // Default password
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, types.Response{Status: "error", Message: err.Error()})
	}

	// Only if the current user role is admin
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, types.Response{Status: "error", Message: "Forbidden"})
	}

	if user.Uuid == "" || user.Name == "" || user.Email == "" || user.Role == "" {
		return c.JSON(http.StatusBadRequest, types.Response{Status: "error", Message: "Uuid, Name, Email, and Role are required"})
	}

	ctx := context.Background()
	// Check if the user exists
	dbUser, err := db.CallQueries.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Response{Status: "error", Message: err.Error()})
	}
	if dbUser.Uuid == "" {
		return c.JSON(http.StatusNotFound, types.Response{Status: "error", Message: "User not found"})
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Response{Status: "error", Message: err.Error()})
	}

	if user.Password == "123456@temp!" {
		hashedPassword = dbUser.Password // Keep the old password
	} else {
		hashedPassword, err = utils.HashPassword(user.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, types.Response{Status: "error", Message: err.Error()})
		}
	}

	// Update the user
	_, err = db.CallQueries.UpdateUser(ctx, db.UpdateUserParams{
		Uuid:      user.Uuid,
		Name:      user.Name,
		Email:     user.Email,
		Password:  hashedPassword,
		Role:      user.Role,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Response{Status: "error", Message: err.Error()})
	}
	responseHTML := `<tr
                        class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
						<td colspan="4" class="px-6 py-4 text-center">
							<div class="flex items-center p-4 mb-4 text-sm text-green-800 rounded-lg bg-green-50 dark:bg-gray-800 dark:text-green-400" role="alert">
								<svg class="flex-shrink-0 inline w-4 h-4 me-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 20 20">
									<path d="M10 .5a9.5 9.5 0 1 0 9.5 9.5A9.51 9.51 0 0 0 10 .5ZM9.5 4a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3ZM12 15H8a1 1 0 0 1 0-2h1v-3H8a1 1 0 0 1 0-2h2a1 1 0 0 1 1 1v4h1a1 1 0 0 1 0 2Z"/>
								</svg>
								<span class="sr-only">Info</span>
								<div>
									<span class="font-medium">Success!</span> User updated successfully. Reloading...
								</div>
							</div>
							<script>
								setTimeout(function() {
									location.reload();
								}, 2000);
							</script>
						</td>
                    </tr>`
	return c.HTML(http.StatusOK, responseHTML)
}
