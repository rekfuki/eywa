package controllers

import (
	"fmt"
	"math"
	"net/http"

	"github.com/davecgh/go-spew/spew"
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

// DeleteFunctions handles function delete requests
func DeleteFunction(c echo.Context) error {
	k8sClient := c.Get("k8s").(*k8s.Client)

	deleteRequest := &types.DeleteFunctionRequest{}
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

// InvocationAlert handles invocation alerts from Alertmanager
func InvocationAlert(c echo.Context) error {
	k8sClient := c.Get("k8s").(*k8s.Client)

	body := types.AlertmanagerAlert{}
	if err := c.Bind(&body); err != nil {
		return err
	}

	for _, alert := range body.Alerts {
		spew.Dump(alert)
		log.WithFields(log.Fields{
			"status":        alert.Status,
			"severity":      alert.Labels.Severity,
			"alert_name":    alert.Labels.AlertName,
			"function_name": alert.Labels.FunctionName,
			"Description":   alert.Annotations.Description,
		}).Info("Received alert")

		fs, err := k8sClient.GetFunctionStatus(alert.Labels.FunctionName)
		if err != nil {
			log.Errorf("Failed to get function status: %s", err)
			return c.JSON(http.StatusInternalServerError, "Internal Server Error")
		}

		if fs == nil {
			return c.JSON(http.StatusNotFound, fmt.Sprintf("Function %q not found", alert.Labels.FunctionName))
		}

		spew.Dump(fs)

		var newReplicas int32
		step := int32(math.Ceil(float64(fs.MaxReplicas) / 100 * float64(fs.ScalingFactor)))
		if body.Status == "firing" && step > 0 {
			newReplicas = fs.Replicas + step
			if newReplicas > fs.MaxReplicas {
				newReplicas = fs.MaxReplicas
			}
		} else {
			newReplicas = fs.MinReplicas
		}

		err = k8sClient.ScaleFunction(fs.Name, newReplicas)
		if err != nil {
			log.Errorf("Failed to scale down: %s", err)
			return c.JSON(http.StatusInternalServerError, "Internal Server Error")
		}
	}

	return c.JSON(http.StatusOK, nil)
}
