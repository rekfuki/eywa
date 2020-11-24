package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"eywa/gateway/clients/k8s"
	"eywa/gateway/types"
)

// CreateFunction deploys a new function deployment to k8s
func CreateFunction(c echo.Context) error {
	k8sClient := c.Get("k8s").(*k8s.Client)

	createRequest := &types.CreateFunctionRequest{}
	if err := c.Bind(createRequest); err != nil {
		return err
	}

	secrets, err := k8sClient.GetSecrets(createRequest.Secrets)
	if err != nil {
		log.Errorf("Failed to get secrets: %s", err)
		// TODO: diff between 500 and 404
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	if err := k8sClient.CreateFunction(createRequest, secrets); err != nil {
		log.Errorf("Failed to deploy function: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusNoContent, nil)
}
