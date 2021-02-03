package controllers

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	ett "eywa/execution-tracker/types"
	"eywa/gateway/clients/k8s"
	"eywa/gateway/types"
	"eywa/go-libs/auth"
	"eywa/go-libs/broker"
	"eywa/go-libs/trigger"
)

// AsyncInvocation dispatches function invocation request
func AsyncInvocation(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)
	bc := c.Get("broker").(*broker.Client)
	functionID := c.Param("function_id")

	filter := k8s.LabelSelector().
		Equals(types.FunctionIDLabel, functionID).
		Equals(types.UserIDLabel, auth.UserID)
	fs, err := k8sClient.GetFunctionStatusFiltered(filter)
	if err != nil {
		log.Errorf("Failed to get functions from k8s: ", err)
		return err
	}

	if fs == nil {
		return c.JSON(http.StatusNotFound, "Function Not Found")
	}

	requestBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	functionName := "UNKNOWN"
	if val, exists := fs.Labels[types.UserDefinedNameLabel]; exists {
		functionName = val
	}

	requestID := c.Request().Header.Get("X-Request-Id")
	payload := broker.QueueRequestMessage{
		Payload: broker.QueueRequest{
			UserID:       auth.UserID,
			RequestID:    requestID,
			Headers:      c.Request().Header,
			Body:         requestBody,
			Path:         c.Param("*"),
			QueryParams:  c.QueryString(),
			FunctionID:   fs.Name,
			FunctionName: functionName,
			CallbackURL:  c.Request().Header.Get("X-Callback-Url"),
			QueuedAt:     time.Now(),
		},
	}

	status := http.StatusAccepted
	eventType := ett.TimelineEventTypeQueued
	eventMessage := types.QueuedMessage(requestID, functionID, functionName)
	if err := bc.ProduceAsync(types.AsyncExecSubject, payload); err != nil {
		log.Errorf("Failed to produce queue request: %s", err)
		status = http.StatusServiceUnavailable
		eventType = ett.TimelineEventTypeFailed
		eventMessage = types.ServerErrorMessage()
	}

	trigger.WithFields(trigger.Fields{
		"function_id": functionID,
		"user_id":     auth.UserID,
		"request_id":  requestID,
		"event_name":  functionName,
		"event_type":  eventType,
		"response":    status,
	}).Fire(types.TimelineHookType)

	trigger.WithFields(trigger.Fields{
		"user_id":       auth.UserID,
		"request_id":    requestID,
		"type":          ett.EventTypeSystem,
		"function_name": functionName,
		"function_id":   functionID,
	}).Fire(types.EventHookType, eventMessage)

	c.Response().Header().Set("X-Request-Id", requestID)
	return c.NoContent(http.StatusAccepted)
}
