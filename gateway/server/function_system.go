package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"eywa/gateway/controllers"
	"eywa/gateway/types"
)

func functionsSystemAPI() []*swagger.Endpoint {
	deployFunction := endpoint.New("POST", "/system/functions", "Deploy function",
		endpoint.Description("Deploy a new function to the cluster"),
		endpoint.Handler(controllers.CreateFunction),
		endpoint.Body(types.CreateFunctionRequest{}, "Deployt function request payload", true),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("Functions"),
	)

	deleteFunction := endpoint.New("DELETE", "/system/functions", "Delete function deployment",
		endpoint.Description("Delete function reployment"),
		endpoint.Handler(controllers.DeleteFunction),
		endpoint.Body(types.DeleteFunctionRequest{}, "Delete function request body", true),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("Functions"),
	)

	scaleFunction := endpoint.New("POST", "/system/functions/scale", "Scale function deployment",
		endpoint.Description("Scale function reployment"),
		endpoint.Handler(controllers.ScaleFunction),
		endpoint.Body(types.ScaleFunctionRequest{}, "Scale function request body", true),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("Functions"),
	)

	return []*swagger.Endpoint{
		deployFunction,
		deleteFunction,
		scaleFunction,
	}
}
