package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"eywa/gateway/clients/k8s"
	"eywa/gateway/types"
)

// DeleteFunctions handles function delete requests
func DeleteFunction(c echo.Context) error {
	k8sClient := c.Get("k8s").(*k8s.Client)

	deleteRequest := &types.DeleteFunctionRequest{}
	if err := c.Bind(deleteRequest); err != nil {
		return err
	}

	function, err := k8sClient.GetFunction(deleteRequest.Name)
	if err != nil {
		log.Errorf("Failed to get function deployment: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	if function == nil {
		return c.JSON(http.StatusNotFound, "Function not found")
	}

	if err = k8sClient.DeleteFunction(function); err != nil {
		log.Errorf("Failed to delete funtion deployment: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusNoContent, nil)
}
