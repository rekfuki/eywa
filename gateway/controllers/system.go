package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"eywa/gateway/clients/k8s"
	"eywa/gateway/types"
)

// SystemDeployFunction deploys a new function deployment to k8s
func SystemDeployFunction(c echo.Context) error {
	k8sClient := c.Get("k8s").(*k8s.Client)

	createRequest := &k8s.DeployFunctionRequest{}
	if err := c.Bind(createRequest); err != nil {
		return err
	}

	secrets, err := k8sClient.GetSecrets(createRequest.Secrets)
	if err != nil {
		log.Errorf("Failed to get secrets: %s", err)
		// TODO: diff between 500 and 404
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	fs, err := k8sClient.DeployFunction(createRequest, secrets)
	if err != nil {
		log.Errorf("Failed to deploy function: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, fs)
}

// SystemDeleteFunction handles function delete requests
func SystemDeleteFunction(c echo.Context) error {
	k8sClient := c.Get("k8s").(*k8s.Client)

	deleteRequest := &types.SystemDeleteFunctionRequest{}
	if err := c.Bind(deleteRequest); err != nil {
		return err
	}

	function, err := k8sClient.GetFunctionStatus(deleteRequest.Name)
	if err != nil {
		log.Errorf("Failed to get function deployment: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	if function == nil {
		return c.JSON(http.StatusNotFound, "Function not found")
	}

	if err = k8sClient.DeleteFunction(deleteRequest.Name); err != nil {
		log.Errorf("Failed to delete funtion deployment: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusNoContent, nil)
}

// SystemScaleFunction scales the function
func SystemScaleFunction(c echo.Context) error {
	k8sClient := c.Get("k8s").(*k8s.Client)

	scaleRequest := &types.SystemScaleFunctionRequest{}
	if err := c.Bind(scaleRequest); err != nil {
		return err
	}

	functionStatus, err := k8sClient.GetFunctionStatus(scaleRequest.Name)
	if err != nil {
		log.Errorf("Failed to get function deployment: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	if functionStatus == nil {
		return c.JSON(http.StatusNotFound, "Function not found")
	}

	if err = k8sClient.ScaleFunction(scaleRequest.Name, scaleRequest.Replicas); err != nil {
		log.Errorf("Failed to delete funtion deployment: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusNoContent, nil)
}
