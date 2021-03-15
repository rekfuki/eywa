package controllers

import (
	"eywa/warden/clients/tugrik"
	"eywa/warden/db"
	"eywa/warden/types"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// ShowLogin shows login page
func ShowLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func loginOK(c echo.Context, user *types.User) error {
	c.Set("user", user)

	_, _, err := setSessionCookie(c)
	if err != nil {
		log.Errorf("Cannot sign token: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusFound, "/")
}

func setSessionCookie(c echo.Context) (string, string, error) {
	sessionSigningKey := c.Get("session_signing_key").(string)
	sessionTimeoutDuration := c.Get("session_timeout_duration").(time.Duration)
	user := c.Get("user").(*types.User)

	id, _ := uuid.NewV4()
	token := &types.Token{
		ID:         id.String(),
		UserID:     user.ID,
		UserName:   user.Name,
		UserAvatar: user.AvatarURL,
		UserEmail:  user.OauthProviderLogin,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(sessionTimeoutDuration).Unix(),
		},
	}

	tokenCookieValue, err := types.SignToken(sessionSigningKey, token)
	if err != nil {
		return "", "", err
	}

	c.SetCookie(&http.Cookie{
		Name:     "eywa-session",
		Value:    tokenCookieValue,
		Domain:   c.Request().Host,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	return token.ID, tokenCookieValue, nil
}

func completeUserAuth(c echo.Context) (*types.User, error) {
	db := c.Get("db").(*db.Client)
	tugrik := c.Get("tugrik").(*tugrik.Client)

	oauthUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return nil, err
	}

	dbUser, err := db.GetUserByOauthUserID(oauthUser.UserID)
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

		if err := tugrik.CreateUserDatabase(id.String()); err != nil {
			// Normally logging should be done at the top level
			// However, this an usual scenario which indicates
			// a system error
			log.Errorf("Failed to create user database: %s", err)

			// Cleanup
			if err := db.DeleteUser(id.String()); err != nil {
				log.Errorf("Failed to cleanup user %s:", err)
				return nil, err
			}

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
