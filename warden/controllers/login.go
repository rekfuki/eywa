package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Show login shows login page
func ShowLogin(c echo.Context) error {
	// claims := &authtypes.Session{
	// 	StandardClaims: jwt.StandardClaims{
	// 		IssuedAt:  time.Now().Unix(),
	// 		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	// 	},
	// }
	// signingKey := c.Get("session-signing-key").(string)
	// passwordLoginEnabled := c.Get("password-login-enabled").(bool)
	// signingKey += "-login-form"
	// token, err := authtoken.SignToken(signingKey, claims)
	// if err != nil {
	// 	// This ought to be impossible
	// 	log.Errorf("Cannot sign token: %s", err)
	// 	return echo.NewHTTPError(http.StatusInternalServerError)
	// }

	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"Token": "fooo",
		"URL":   c.FormValue("url"),
	})
}
