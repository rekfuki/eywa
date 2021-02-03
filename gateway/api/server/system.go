package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"eywa/gateway/api/controllers"
	"eywa/gateway/types"
)

func systemAPI() []*swagger.Endpoint {
	getFunctions := endpoint.New("GET", "/system/functions", "Get function",
		endpoint.Description("Get function"),
		endpoint.Handler(controllers.SystemGetFunctions),
		endpoint.Response(http.StatusOK, types.MultiFunctionStatusResponse{}, "Success"),
		endpoint.Tags("System"),
	)

	scaleFunction := endpoint.New("POST", "/system/functions/{function_id}/scale/{replicas}", "Scale function deployment",
		endpoint.Description("Scale function reployment"),
		endpoint.Handler(controllers.SystemScaleFunction),
		endpoint.PathMap(map[string]swagger.Parameter{
			"function_id": {
				Type:   "string",
				Format: "uuid",
			},
			"replicas": {
				Type:    "integer",
				Minimum: &[]int64{0}[0],
			},
		}),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("System"),
	)

	return []*swagger.Endpoint{
		getFunctions,
		scaleFunction,
	}
}
