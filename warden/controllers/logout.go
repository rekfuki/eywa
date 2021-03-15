package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Logout logs the user out
func Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "eywa-session",
		Domain:   c.Request().Host,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	})

	return c.NoContent(http.StatusOK)
}
