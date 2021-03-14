package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"eywa/gateway/metrics"
	"eywa/gateway/types"
	"eywa/go-libs/auth"
)

// QueryMetrics returns prometheus metrics
func QueryMetrics(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	mc := c.Get("metrics").(*metrics.Client)

	var qr types.QueryRequest
	if err := c.Bind(&qr); err != nil {
		return err
	}

	validationErrors := map[string][]string{}
	if qr.Type == "range" {
		if qr.Start == 0 {
			validationErrors["start"] = append(validationErrors["star"], "is required for type range queries")
		}

		if qr.End == 0 {
			validationErrors["end"] = append(validationErrors["star"], "is required for type range queries")
		}

		if qr.Start == 0 {
			validationErrors["start"] = append(validationErrors["star"], "must be before end")
		}

		if qr.Step <= 0 {
			validationErrors["step"] = append(validationErrors["step"], "must be a positive integer")
		}
	}

	if len(validationErrors) > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation error",
			"details": validationErrors,
		})
	}

	if strings.Contains(qr.LabelMatchers, "user_id") {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation error",
			"details": map[string]string{
				"label_matchers": "user_id not allowed as label matcher",
			},
		})
	}

	// Scope to user only
	qr.LabelMatchers += fmt.Sprintf(`,user_id="%s"`, auth.UserID)

	tmpl, err := template.New("test").Delims("<<", ">>").Parse(qr.Query)
	if err != nil {
		log.Errorf("Failed to parse metrics template: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, qr); err != nil {
		log.Errorf("Failed to buuild metrics template: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	tmplString := tpl.String()

	var query string
	var endpoint string
	if qr.Type == "range" {
		endpoint = "/query_range"
		query = fmt.Sprintf("query=%s&start=%f&end=%f&step=%d", tmplString, qr.Start, qr.End, qr.Step)
	} else {
		endpoint = "/query"
		query = fmt.Sprintf("query=%s", tmplString)
		if qr.End != 0 {
			query += fmt.Sprintf("&time=%f", qr.End)
		}
	}

	fmt.Println(query)
	result, qErr := mc.QueryMetrics(endpoint, query)
	if qErr != nil {
		if qErr.Type == types.ErrTypeBadRequest {
			return echo.NewHTTPError(http.StatusBadRequest, qErr.String())
		}

		log.Errorf("Failed to get metrics from prometheus: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, result)
}
