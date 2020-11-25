package controllers

import (
	"eywa/gateway/clients/k8s"
	"eywa/gateway/types"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func ScaleFunction(c echo.Context) error {
	k8sClient := c.Get("k8s").(*k8s.Client)

	scaleRequest := &types.ScaleFunctionRequest{}
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
