package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"eywa/gateway/controllers"
	"eywa/gateway/types"
)

func functionsAPI() []*swagger.Endpoint {
	getFunctions := endpoint.New("GET", "/functions", "Get functions",
		endpoint.Description("Get all function belonging to a user"),
		endpoint.Handler(controllers.GetFunctions),
		endpoint.Response(http.StatusOK, types.MultiFunctionStatusResponse{}, "Success"),
		endpoint.Tags("Functions"),
	)

	getFunction := endpoint.New("GET", "/functions/{function_id}", "Get specific function",
		endpoint.Description("Get a function belonging to a user"),
		endpoint.Handler(controllers.GetFunction),
		endpoint.Path("function_id", "string", "uuid", "UUID of a function"),
		endpoint.Response(http.StatusOK, types.FunctionStatusResponse{}, "Success"),
		endpoint.Tags("Functions"),
	)

	deployFunction := endpoint.New("POST", "/functions", "Create a function",
		endpoint.Description("Create a function from an image"),
		endpoint.Handler(controllers.DeployFunction),
		endpoint.Body(types.DeployFunctionRequest{}, "Function deployment payload", true),
		endpoint.Response(http.StatusCreated, types.FunctionStatusResponse{}, "Success"),
		endpoint.Tags("Functions"),
	)

	updateFunction := endpoint.New("PUT", "/functions/{function_id}", "Update a function",
		endpoint.Description("Update a deployed function from k8s"),
		endpoint.Handler(controllers.UpdateFunction),
		endpoint.Path("function_id", "string", "uuid", "UUID of a function"),
		endpoint.Body(types.UpdateFunctionRequest{}, "Function update payload", true),
		endpoint.Response(http.StatusOK, types.FunctionStatusResponse{}, "Success"),
		endpoint.Tags("Functions"),
	)

	deleteFunction := endpoint.New("DELETE", "/functions/{function_id}", "Delete a function",
		endpoint.Description("Delete a deployed function from k8s"),
		endpoint.Handler(controllers.DeleteFunction),
		endpoint.Path("function_id", "string", "uuid", "UUID of a function"),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("Functions"),
	)

	return []*swagger.Endpoint{
		getFunctions,
		getFunction,
		deployFunction,
		updateFunction,
		deleteFunction,
	}
}
