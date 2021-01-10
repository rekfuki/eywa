package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"eywa/gateway/controllers"
	"eywa/gateway/types"
)

func functionsAPI() []*swagger.Endpoint {
	deployFunction := endpoint.New("POST", "/functions", "Create a function",
		endpoint.Description("Create a function from an image"),
		endpoint.Handler(controllers.DeployFunction),
		endpoint.Body(types.DeployFunctionRequest{}, "Service deployment payload", true),
		endpoint.Response(http.StatusCreated, types.FunctionStatusResponse{}, "Success"),
		endpoint.Tags("Functions"),
	)

	// getImage := endpoint.New("GET", "/registry/images/{image_id}", "Get a specific image",
	// 	endpoint.Description("Get a specific image"),
	// 	endpoint.Handler(controllers.GetImage),
	// 	endpoint.Path("image_id", "string", "uuid", "UUID of an image"),
	// 	endpoint.Response(http.StatusOK, types.Image{}, "Success"),
	// 	endpoint.Tags("Registry"),
	// )

	// createImage := endpoint.New("POST", "/registry/images", "Create an image",
	// 	endpoint.Description("Create an image"),
	// 	endpoint.Consumes("multipart/form-data"),
	// 	endpoint.Handler(controllers.CreateImage),
	// 	endpoint.FormDataMap(map[string]swagger.Parameter{
	// 		"source": {
	// 			Type:        "string",
	// 			Format:      "file",
	// 			Description: "Zip file containing source",
	// 			Required:    true,
	// 		},
	// 		"version": {
	// 			Type:        "string",
	// 			Format:      "string",
	// 			Pattern:     "^(\\d{1,3}\\.?){3}$",
	// 			Description: "Source code semantic version",
	// 			Required:    true,
	// 		},
	// 		"language": {
	// 			Type:        "string",
	// 			Format:      "string",
	// 			Enum:        []string{"Go"},
	// 			Description: "Language the source is written in",
	// 			Required:    true,
	// 		},
	// 		"name": {
	// 			Type:        "string",
	// 			Format:      "string",
	// 			Description: "Name of the service",
	// 			MinLength:   5,
	// 			Required:    true,
	// 		},
	// 	}),
	// 	endpoint.Response(http.StatusNoContent, types.Image{}, "Success"),
	// 	endpoint.Tags("Registry"),
	// )

	// deleteImage := endpoint.New("DELETE", "/registry/images/{image_id}", "Delete an image",
	// 	endpoint.Description("Delete an image"),
	// 	endpoint.Handler(controllers.DeleteImage),
	// 	endpoint.Path("image_id", "string", "uuid", "UUID of an image"),
	// 	endpoint.Response(http.StatusNoContent, "", "Success"),
	// 	endpoint.Tags("Registry"),
	// )

	return []*swagger.Endpoint{
		deployFunction,
	}
}
