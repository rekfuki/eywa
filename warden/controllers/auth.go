package controllers

import (
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"eywa/warden/db"
	"eywa/warden/types"
)

// OAuth handles oauth requests
func OAuth(c echo.Context) error {
	// try to get the user without re-authenticating
	user, err := completeUserAuth(c)
	if err != nil {
		gothic.BeginAuthHandler(c.Response(), c.Request())
	}

	spew.Dump(user)
	return loginOK(c, user)
}

// CompleteOAuth completes the oath chain
func CompleteOAuth(c echo.Context) error {
	user, err := completeUserAuth(c)
	if err != nil {
		log.Errorf("Failed to complete user auth: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	spew.Dump(user)
	return loginOK(c, user)
}

func loginOK(c echo.Context, user *types.User) error {
	c.Set("user", user)
	_, _, err := setSessionCookie(c)
	if err != nil {
		log.Errorf("Cannot sign token: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// TODO: Redirect to previous location if existed
	return c.Redirect(http.StatusFound, "/dashboard")
}

func setSessionCookie(c echo.Context) (string, string, error) {
	sessionSigningKey := c.Get("session_signing_key").(string)
	sessionTimeoutDuration := c.Get("session_timeout_duration").(time.Duration)
	user := c.Get("user").(*types.User)

	id, _ := uuid.NewV4()
	session := &types.Session{
		ID:     id.String(),
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * sessionTimeoutDuration).Unix(),
		},
	}

	tokenCookieValue, err := types.SignToken(sessionSigningKey, session)
	if err != nil {
		return "", "", err
	}

	c.SetCookie(&http.Cookie{
		Name:     "eywa-token",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	})

	c.SetCookie(&http.Cookie{
		Name:     "eywa-token",
		Value:    tokenCookieValue,
		Domain:   c.Request().Host,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	return session.ID, tokenCookieValue, nil
}

func completeUserAuth(c echo.Context) (*types.User, error) {
	db := c.Get("db").(*db.Client)

	oauthUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return nil, err
	}
	spew.Dump(oauthUser)

	dbUser, err := db.GetUserByUserID(oauthUser.UserID)
	if err != nil {
		return nil, err
	}

	if dbUser == nil {
		now := time.Now()
		id, _ := uuid.NewV4()
		newUser := &types.User{
			ID:                 id.String(),
			Name:               oauthUser.Name,
			AvatarURL:          oauthUser.AvatarURL,
			OauthProvider:      oauthUser.Provider,
			OauthProviderID:    oauthUser.UserID,
			OauthProviderEmail: oauthUser.Email,
			OauthProviderLogin: oauthUser.NickName,
			CreatedAt:          now,
			LastSeenAt:         now,
		}

		dbUser, err = db.CreateUser(newUser)
		if err != nil {
			return nil, err
		}
	} else {
		dbUser = &types.User{
			ID:                 dbUser.ID,
			Name:               oauthUser.Name,
			AvatarURL:          oauthUser.AvatarURL,
			OauthProvider:      oauthUser.Provider,
			OauthProviderID:    oauthUser.UserID,
			OauthProviderEmail: oauthUser.Email,
			OauthProviderLogin: oauthUser.NickName,
			CreatedAt:          dbUser.CreatedAt,
			LastSeenAt:         time.Now(),
		}

		if err = db.UpdateUser(dbUser); err != nil {
			log.Errorf("Failed to update user: %s", err)
			return nil, err
		}
	}

	return dbUser, nil
}
