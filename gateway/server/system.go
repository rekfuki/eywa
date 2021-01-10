package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"eywa/gateway/clients/k8s"
	"eywa/gateway/controllers"
	"eywa/gateway/types"
)

func systemAPI() []*swagger.Endpoint {
	deployFunction := endpoint.New("POST", "/system/functions", "Deploy function",
		endpoint.Description("Deploy a new function to the cluster"),
		endpoint.Handler(controllers.SystemDeployFunction),
		endpoint.Body(k8s.DeployFunctionRequest{}, "Deployt function request payload", true),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("System"),
	)

	deleteFunction := endpoint.New("DELETE", "/system/functions", "Delete function deployment",
		endpoint.Description("Delete function reployment"),
		endpoint.Handler(controllers.SystemDeleteFunction),
		endpoint.Body(types.SystemDeleteFunctionRequest{}, "Delete function request body", true),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("System"),
	)

	scaleFunction := endpoint.New("POST", "/system/functions/scale", "Scale function deployment",
		endpoint.Description("Scale function reployment"),
		endpoint.Handler(controllers.SystemScaleFunction),
		endpoint.Body(types.SystemScaleFunctionRequest{}, "Scale function request body", true),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("System"),
	)

	// getResourceLimits := endpoint.New("GET", "/system/functions/limits", "Get function resource limits",
	// 	endpoint.Description("Get function resource limits"),
	// 	endpoint.Handler(controllers.SystemGetResourceLimits),
	// 	endpoint.Response(http.StatusOK, types.ResourceLimits{}, "Success"),
	// 	endpoint.Tags("Functions"),
	// )

	// getFunctionsFiltered := endpoint.New("POST", "/system/functions/query", "Get filtered function",
	// 	endpoint.Description("Get filetered function"),
	// 	endpoint.Handler(controllers.SystemGetFunctionsFiltered),
	// 	endpoint.Body(types.GetFunctionsFilteredRequest{}, "Payload to filter by", true),
	// 	endpoint.Response(http.StatusOK, types.GetFunctionsFilteredResponse{}, "Success"),
	// 	endpoint.Tags("Functions"),
	// )

	return []*swagger.Endpoint{
		deployFunction,
		deleteFunction,
		scaleFunction,
		// getResourceLimits,
		// getFunctionsFiltered,
	}
}
