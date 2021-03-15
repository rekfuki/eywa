package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	gt "eywa/gateway/types"
	"eywa/go-libs/auth"
	"eywa/tugrik/clients/gateway"
	"eywa/tugrik/db"
	"eywa/tugrik/types"
)

// CreateUserDatabase creates a user and their database in mongo
func CreateUserDatabase(c echo.Context) error {
	db := c.Get("db").(*db.Client)
	gc := c.Get("gateway").(*gateway.Client)
	userID := c.Param("user_id")
	password, err := uuid.NewV4()
	if err != nil {
		log.Errorf("Failed to generate uuid: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	exists, err := db.CheckUserExists(userID)
	if err != nil {
		log.Errorf("Failed to check user existence: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	if exists {
		return c.JSON(http.StatusBadRequest, "User database already exists")
	}

	if err := db.CreateUserDatabase(userID, password.String()); err != nil {
		log.Errorf("Failed to create user database: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	secretRequest := gt.CreateSecretRequest{
		Name: "mongodb-credentials",
		Data: map[string]string{
			"username": userID,
			"password": password.String(),
			"database": userID,
		},
	}

	if err := gc.CreateDatabaseSecret(userID, secretRequest); err != nil {
		log.Errorf("Failed to create database secret: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, types.CreateUserDatabaseResponse{
		Username: userID,
		Database: userID,
		Password: password.String(),
	})
}

// GetUserDatabase returns stats about users database
func GetUserDatabase(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	db := c.Get("db").(*db.Client)

	exists, err := db.CheckUserExists(auth.UserID)
	if err != nil {
		log.Errorf("Failed to check user existence: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	if !exists {
		return c.JSON(http.StatusNotFound, "Database Not Found")
	}

	stats, err := db.GetUserDatabase(auth.UserID)
	if err != nil {
		log.Errorf("Failed to get user db: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	stats.UserID = auth.UserID
	return c.JSON(http.StatusOK, stats)
}

// DeleteCollection deletes a collection from user database
func DeleteCollection(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	db := c.Get("db").(*db.Client)
	collectionName := c.Param("collection_name")

	exists, err := db.CheckUserExists(auth.UserID)
	if err != nil {
		log.Errorf("Failed to check user existence: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	if !exists {
		return c.JSON(http.StatusNotFound, "Database Not Found")
	}

	if err := db.DropCollection(auth.UserID, collectionName); err != nil {
		log.Errorf("Failed to get user db: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.NoContent(http.StatusNoContent)
}
