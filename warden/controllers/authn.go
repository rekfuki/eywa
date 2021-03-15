package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"eywa/warden/authn"
	"eywa/warden/db"
	"eywa/warden/types"
)

const timeFormat = time.RFC3339Nano

// CheckAuth handles the authn request from ingress
func CheckAuth(c echo.Context) error {
	db := c.Get("db").(*db.Client)
	signingKey := c.Get("session_signing_key").(string)
	sessionTimeoutduration := c.Get("session_timeout_duration").(time.Duration)

	authType, token, authErr := authn.IsReqAuthenticated(c.Request(), signingKey)
	if authErr != nil {
		if authErr.Type == types.ErrTypeAuthenticationError {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	c.Response().Header().Add("X-Eywa-User-Id", token.UserID)
	c.Response().Header().Add("X-Eywa-User-Name", token.UserName)
	c.Response().Header().Add("X-Eywa-User-Email", token.UserEmail)
	c.Response().Header().Add("X-Eywa-User-Avatar", token.UserAvatar)

	if authType == authn.AuthTypeSession {
		now := time.Now().UTC().Unix()
		if now-token.IssuedAt > (token.ExpiresAt-now)*2 {
			token.IssuedAt = now
			token.ExpiresAt = now + int64(sessionTimeoutduration.Seconds())*60
		}
	}

	tokenStr, err := types.SignToken(signingKey, token)
	if err != nil {
		log.Errorf("Failed to sign token: %s", err)
		echo.NewHTTPError(http.StatusInternalServerError)
	}

	if authType == authn.AuthTypeSession {
		c.SetCookie(&http.Cookie{
			Name:     "eywa-session",
			Value:    tokenStr,
			Domain:   c.Request().Host,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		})
	} else if authType == authn.AuthTypeToken {
		c.Response().Header().Add("X-Eywa-Token", tokenStr)
	}

	go func() {
	}()

	c.Response().After(func() {
		if err := db.SetUserLastSeenAt(token.UserID, time.Now()); err != nil {
			log.Errorf("Failed to update user last_seen_at: %s", err)
		}
	})

	return c.NoContent(http.StatusOK)
}
