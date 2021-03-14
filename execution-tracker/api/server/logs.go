package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"eywa/execution-tracker/api/controllers"
	"eywa/execution-tracker/types"
)

func logsAPI() []*swagger.Endpoint {
	timelineList := endpoint.New("GET", "/timeline", "Get timeline list",
		endpoint.Handler(controllers.GetInvocationList),
		endpoint.QueryMap(map[string]swagger.Parameter{
			"page": {
				Type:        "integer",
				Minimum:     &[]int64{1}[0],
				Description: "Page number to return",
			},
			"per_page": {
				Type:        "integer",
				Minimum:     &[]int64{0}[0],
				Description: "Number of records per page",
			},
			"function_id": {
				Type:        "string",
				Format:      "uuid",
				Description: "UUID of a function",
			},
			"query": {
				Type:        "string",
				Description: "Query string to search by",
			},
		}),
		endpoint.Response(http.StatusOK, types.TimelineLogsResponse{}, "success"),
		endpoint.Tags("Timeline"),
	)

	timelineDetails := endpoint.New("GET", "/timeline/{request_id}", "Get details of a timeline",
		endpoint.Handler(controllers.GetInvocation),
		endpoint.Path("request_id", "string", "uuid", "UUID of a request"),
		endpoint.Response(http.StatusOK, types.TimelineDetails{}, "success"),
		endpoint.Tags("Timeline"),
	)

	eventLogsQuery := endpoint.New("POST", "/events/query", "Query event log records",
		endpoint.Handler(controllers.GetEventLogs),
		endpoint.QueryMap(map[string]swagger.Parameter{
			"page": {
				Type:        "integer",
				Minimum:     &[]int64{1}[0],
				Description: "Page number to return",
			},
			"per_page": {
				Type:        "integer",
				Minimum:     &[]int64{0}[0],
				Description: "Number of records per page",
			},
		}),
		endpoint.Body(types.EventLogsQuery{}, "Query object that needs to be added", false),
		endpoint.Response(http.StatusOK, types.EventLogsResponse{}, "success"),
		endpoint.Tags("Events"),
	)

	return []*swagger.Endpoint{
		timelineList,
		timelineDetails,
		eventLogsQuery,
	}
}
