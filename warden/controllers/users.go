package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"eywa/warden/db"
)

// GetUser returns user
func GetUser(c echo.Context) error {
	db := c.Get("db").(*db.Client)
	userID := c.Request().Header.Get("X-Eywa-User-Id")

	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	user, err := db.GetUserByInternalUserID(userID)
	if err != nil {
		log.Errorf("Failed to get user from db: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User Not Found")
	}

	return c.JSON(http.StatusOK, user)
}
