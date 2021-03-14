package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"

	ett "eywa/execution-tracker/types"
	"eywa/gateway/clients/k8s"
	"eywa/gateway/metrics"
	"eywa/gateway/types"
	"eywa/go-libs/auth"
	"eywa/go-libs/trigger"
	wet "eywa/watchdog/executor"
)

// Proxy proxies requests from outside to the corret internal pod
func Proxy(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodGet:

		return proxyRequest(c)
	default:
		return c.JSON(http.StatusMethodNotAllowed, nil)
	}
}

func proxyRequest(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	proxyClient := c.Get("proxy").(*resty.Client)
	k8s := c.Get("k8s").(*k8s.Client)
	metrics := c.Get("metrics").(*metrics.Client)
	functionName := c.Get("function_name").(string)

	functionID := c.Param("function_id")
	if functionID == "" {
		return c.JSON(http.StatusBadRequest, "Missing function id")
	}

	requestID := c.Request().Header.Get("X-Request-ID")
	defaultEventFields := trigger.Fields{
		"user_id":       auth.UserID,
		"request_id":    requestID,
		"type":          ett.EventTypeSystem,
		"function_name": functionName,
		"function_id":   functionID,
	}

	defaultTimelineFields := trigger.Fields{
		"user_id":     auth.UserID,
		"request_id":  requestID,
		"function_id": functionID,
		"method":      c.Request().Method,
	}

	requestBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	stripHeaders(c.Request().Header)
	path := "/" + c.Param("*")
	trigger.WithFields(defaultEventFields).WithFields(trigger.Fields{
		"path":         path,
		"query_params": c.QueryString(),
		"body":         requestBody,
		"headers":      c.Request().Header,
		"message":      types.SyncExecutionStartMessage(functionName),
	}).Fire(types.EventHookType)

	metrics.ObserveInvocationStarted(functionID, functionName,
		auth.UserID, path, c.Request().Method)

	fullChainStart := time.Now()

	functionAddr, err := k8s.Resolve(functionID)
	if err != nil {
		log.Errorf("k8s error: cannot find %s: %s\n", functionID, err)

		trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
			"event_name": functionName,
			"event_type": ett.TimelineEventTypeSystemError,
			"response":   http.StatusServiceUnavailable,
			"duration":   time.Since(fullChainStart).Milliseconds(),
		}).Fire(types.TimelineHookType)

		trigger.WithFields(defaultEventFields).WithFields(trigger.Fields{
			"is_error": true,
			"message":  types.ServerErrorMessage(),
		}).Fire(types.EventHookType)

		return echo.NewHTTPError(http.StatusServiceUnavailable)
	}

	url := fmt.Sprintf("%s%s", functionAddr, path)
	proxyRequest := proxyClient.R().SetQueryString(c.QueryString())

	if len(requestBody) > 0 {
		proxyRequest.SetBody(requestBody)
	}

	copyHeaders(proxyRequest.Header, &c.Request().Header)

	proxyStart := time.Now()
	var result wet.FunctionResponse
	response, err := proxyRequest.
		SetResult(&result).
		Execute(c.Request().Method, url)
	if err != nil || response.IsError() {
		log.Errorf("Error with proxy request to: %s, %s\n", proxyRequest.URL, err)

		trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
			"event_name": functionName,
			"event_type": ett.TimelineEventTypeSystemError,
			"response":   http.StatusServiceUnavailable,
			"duration":   time.Since(fullChainStart).Milliseconds(),
		}).Fire(types.TimelineHookType)

		trigger.WithFields(defaultEventFields).WithFields(trigger.Fields{
			"is_error": true,
			"message":  types.ServerErrorMessage(),
		}).Fire(types.EventHookType)

		return c.JSON(http.StatusServiceUnavailable, "Service Unavailable")
	}

	proxyFinish := time.Since(proxyStart)
	log.Infof("%s took %f seconds\n", functionID, proxyFinish.Seconds())

	metrics.ObserveInvocationComplete(functionID, functionName, auth.UserID, path, result.Status, proxyFinish)

	eventType := ett.TimelineEventTypeFinished
	if result.Status >= 400 {
		eventType = ett.TimelineEventTypeFailed
	}

	trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
		"event_name": functionName,
		"event_type": eventType,
		"response":   result.Status,
		"duration":   proxyFinish.Milliseconds(),
	}).Fire(types.TimelineHookType)

	defaultEventFields["type"] = ett.EventTypeUser
	defaultEventFields["created_at"] = fullChainStart.Add(proxyFinish)

	// Inject back request id header
	if result.Headers == nil {
		result.Headers = http.Header{}
	}

	result.Headers.Set("X-Request-Id", requestID)
	result.Headers.Del("X-Eywa-Token")

	trigger.WithFields(defaultEventFields).WithFields(trigger.Fields{
		"is_error": eventType == ett.TimelineEventTypeFailed,
		"status":   result.Status,
		"headers":  result.Headers,
		"body":     result.Body,
		"stdout":   result.Stdout,
		"stderr":   result.Stderr,
		"message":  types.SyncExecutionFinishMessage(functionName, result.Status, proxyFinish),
	}).Fire(types.EventHookType)

	return copyResponse(c, result.Status, result.Headers, result.Body)
}

func copyResponse(c echo.Context, statusCode int, headers http.Header, body []byte) error {
	h := c.Response().Header()
	for k, v := range headers {
		h[k] = v
	}
	return c.Blob(statusCode, headers.Get("Content-Type"), body)
}
