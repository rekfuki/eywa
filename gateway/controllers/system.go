package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"eywa/gateway/clients/k8s"
	"eywa/gateway/types"
	"eywa/go-libs/auth"
)

// SystemGetFunctions returns list of functions
func SystemGetFunctions(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)

	if !auth.IsOperator() {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}
	fss, err := k8sClient.GetFunctionsStatus()
	if err != nil {
		log.Errorf("Failed to get functions from k8s: ", err)
		return err
	}

	sfss := []types.FunctionStatusResponse{}
	for _, fs := range fss {
		sfss = append(sfss, makeFunctionStatusResponse(&fs))
	}

	return c.JSON(http.StatusOK, types.MultiFunctionStatusResponse{
		Objects: sfss,
		Total:   len(sfss),
	})
}

// SystemScaleFunction scales the function
func SystemScaleFunction(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)
	functionID := c.Param("function_id")
	replicasStr := c.Param("replicas")

	if !auth.IsOperator() {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	replicas, err := strconv.ParseInt(replicasStr, 10, 64)
	if err != nil {
		log.Errorf("Failed to parse replica count %s: %s", replicasStr, err)
		return err
	}

	filter := k8s.LabelSelector().Equals(types.FunctionIDLabel, functionID)
	functionStatus, err := k8sClient.GetFunctionStatus(filter)
	if err != nil {
		log.Errorf("Failed to get function deployment: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	if functionStatus == nil {
		return c.JSON(http.StatusNotFound, "Function not found")
	}

	if err = k8sClient.ScaleFunction(filter, int(replicas)); err != nil {
		log.Errorf("Failed to scale funtion deployment: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.NoContent(http.StatusNoContent)
}
