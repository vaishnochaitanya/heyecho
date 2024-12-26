package handlers

import (
	"github.com/labstack/echo/v4"
	"vchan.in/heyecho/views/foo"
)

var profile = "guest"

func HandleFoo(c echo.Context) error {
	return render(c, foo.Foo(profile))
}
