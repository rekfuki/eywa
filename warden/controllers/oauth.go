package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	log "github.com/sirupsen/logrus"
)

// OAuth handles oauth requests
func OAuth(c echo.Context) error {
	if c.QueryParam("code") != "" {
		return completeOAuth(c)
	}

	// try to get the user without re-authenticating
	user, err := completeUserAuth(c)
	if err == nil {
		return loginOK(c, user)
	}

	gothic.BeginAuthHandler(c.Response(), c.Request())
	return nil
}

// CompleteOAuth completes the oath chain
func completeOAuth(c echo.Context) error {
	user, err := completeUserAuth(c)
	if err != nil {
		log.Errorf("Failed to complete user auth: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return loginOK(c, user)
}
