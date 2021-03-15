package controllers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"

	"eywa/go-libs/auth"
	"eywa/warden/db"
	"eywa/warden/types"
)

// GetAccessTokens returns all user access tokens
func GetAccessTokens(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	db := c.Get("db").(*db.Client)
	perPage := c.Get("per_page").(int)
	pageNumber := c.Get("page_number").(int)
	query := c.QueryParam("query")

	total, accessTokens, err := db.GetAccessTokens(auth.UserID, query, pageNumber, perPage)
	if err != nil {
		log.Errorf("Failed to get access tokens from db: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, types.AccessTokensResponse{
		Objects: accessTokens,
		Total:   total,
		PerPage: perPage,
		Page:    pageNumber,
	})
}

// CreateToken creates a new access token to be used with the API
func CreateToken(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	db := c.Get("db").(*db.Client)
	sessionSigningKey := c.Get("session_signing_key").(string)

	var cr types.CreateTokenRequest
	if err := c.Bind(&cr); err != nil {
		return err
	}

	now := time.Now()

	if cr.ExpiresAt > 0 && cr.ExpiresAt < now.Unix() {
		return c.JSON(http.StatusBadRequest, "Expiry time cannot be in the past")
	}

	id, _ := uuid.NewV4()
	token := &types.Token{
		ID:     id.String(),
		UserID: auth.UserID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: now.Unix(),
		},
	}

	if cr.ExpiresAt > 0 {
		token.StandardClaims.ExpiresAt = cr.ExpiresAt
	}

	tokenString, err := types.SignToken(sessionSigningKey, token)
	if err != nil {
		log.Errorf("Failed to sign access token: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	nameID, _ := uuid.NewV4()
	accessToken := &types.AccessToken{
		ID:        nameID.String(),
		UserID:    auth.UserID,
		Name:      nameID.String()[:9] + cr.Name,
		Token:     tokenString,
		CreatedAt: now.Unix(),
		ExpiresAt: cr.ExpiresAt,
	}

	if err := db.CreateAccessToken(accessToken); err != nil {
		log.Errorf("Failed to create access token: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, accessToken)
}

// DeleteAccessToken deletes user access token
func DeleteAccessToken(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	db := c.Get("db").(*db.Client)
	tokenID := c.Param("token_id")

	accessToken, err := db.GetAccessToken(auth.UserID, tokenID)
	if err != nil {
		log.Errorf("Failed to get access token from db: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	if accessToken == nil {
		return c.JSON(http.StatusNotFound, "Access Token Not Found")
	}

	if err := db.DeleteAccessToken(auth.UserID, tokenID); err != nil {
		log.Errorf("Failed to delete access token from db: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.NoContent(http.StatusNoContent)
}
