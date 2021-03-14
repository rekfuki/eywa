package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"eywa/gateway/api/controllers"
	"eywa/gateway/types"
)

func metricsAPI() []*swagger.Endpoint {
	getMetrics := endpoint.New("POST", "/metrics/query", "Query metrics",
		endpoint.Description("Query prometheus metrics"),
		endpoint.Handler(controllers.QueryMetrics),
		endpoint.Body(types.QueryRequest{}, "Query payload", true),
		endpoint.Response(http.StatusOK, types.QueryResponse{}, "Success"),
		endpoint.Tags("Metrics"),
	)

	return []*swagger.Endpoint{
		getMetrics,
	}
}
