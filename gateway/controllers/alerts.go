package controllers

import (
	"fmt"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"eywa/gateway/clients/k8s"
	"eywa/gateway/types"
)

// InvocationAlert handles invocation alerts from Alertmanager
func InvocationAlert(c echo.Context) error {
	k8sClient := c.Get("k8s").(*k8s.Client)

	body := types.AlertmanagerAlert{}
	if err := c.Bind(&body); err != nil {
		return err
	}

	for _, alert := range body.Alerts {
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

		var newReplicas int
		step := int(math.Ceil(float64(fs.MaxReplicas) / 100 * float64(fs.ScalingFactor)))
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
