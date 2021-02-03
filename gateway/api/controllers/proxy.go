package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

	path := c.Param("*")
	trigger.WithFields(defaultEventFields).FireForEach(types.EventHookType,
		types.StartMessage(requestID, functionID, functionName),
		types.RequestContextMessage(path, c.QueryString(), requestBody, c.Request().Header),
	)

	metrics.Observe(c.Request().Method, functionName, auth.UserID,
		http.StatusProcessing, "started", time.Second*0)

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

		trigger.WithFields(defaultEventFields).
			WithFields(trigger.Fields{"is_error": true}).
			Fire(types.EventHookType, types.ServerErrorMessage())

		return c.JSON(http.StatusServiceUnavailable, err)
	}

	url := fmt.Sprintf("%s/%s", functionAddr, path)
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
	if err != nil {
		log.Errorf("Error with proxy request to: %s, %s\n", proxyRequest.URL, err)

		trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
			"event_name": functionName,
			"event_type": ett.TimelineEventTypeSystemError,
			"response":   http.StatusServiceUnavailable,
			"duration":   time.Since(fullChainStart).Milliseconds(),
		}).Fire(types.TimelineHookType)

		trigger.WithFields(defaultEventFields).
			WithFields(trigger.Fields{"is_error": true}).
			Fire(types.EventHookType, types.ServerErrorMessage())

		return c.JSON(http.StatusServiceUnavailable, "Service Unavailable")
	}

	proxyFinish := time.Since(proxyStart)
	log.Infof("%s took %f seconds\n", functionID, proxyFinish.Seconds())

	metrics.Observe(c.Request().Method, functionName, auth.UserID,
		http.StatusOK, "completed", time.Since(fullChainStart))

	eventType := ett.TimelineEventTypeFinished
	if response.IsError() {
		eventType = ett.TimelineEventTypeFailed
	}

	trigger.WithFields(defaultTimelineFields).WithFields(trigger.Fields{
		"event_name": functionName,
		"event_type": eventType,
		"response":   response.StatusCode(),
		"duration":   proxyFinish.Milliseconds(),
	}).Fire(types.TimelineHookType)

	defaultEventFields["type"] = ett.EventTypeUser
	defaultEventFields["created_at"] = fullChainStart.Add(proxyFinish)

	fields := []interface{}{
		types.ResponseContextMessage(path, c.QueryString(),
			result.Body, response.Header(), response.StatusCode()),
	}

	if len(result.Stdout) > 0 {
		message := strings.Join(result.Stdout, "\n")
		fields = append(fields, types.StdoutMessage(message))
	}

	if len(result.Stderr) > 0 {
		message := strings.Join(result.Stderr, "\n")
		fields = append(fields, types.StderrMessage(message))
	}

	trigger.WithFields(defaultEventFields).
		WithFields(trigger.Fields{"is_error": eventType == ett.TimelineEventTypeFailed}).
		FireForEach(types.EventHookType, fields...)

	response.Header().Del("Content-Length")
	return copyResponse(c, response, result.Body)
}

func copyHeaders(destination http.Header, source *http.Header) {
	for k, v := range *source {
		vClone := make([]string, len(v))
		copy(vClone, v)
		(destination[k]) = v
	}
}

func copyResponse(c echo.Context, response *resty.Response, body []byte) error {
	h := c.Response().Header()
	for k, v := range response.Header() {
		h[k] = v
	}
	return c.Blob(response.StatusCode(), response.Header().Get("Content-Type"), body)
}
