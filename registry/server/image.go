package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"eywa/registry/controllers"
	"eywa/registry/types"
)

func imagesAPI() []*swagger.Endpoint {
	getImages := endpoint.New("GET", "/images", "Get all the images",
		endpoint.Description("Get all the images"),
		endpoint.Handler(controllers.GetImages),
		endpoint.QueryMap(map[string]swagger.Parameter{
			"page": {
				Type:        "integer",
				Description: "Page number to return",
			},
			"per_page": {
				Type:        "integer",
				Description: "Number of records per page",
			},
		}),
		endpoint.Response(http.StatusOK, types.GetImagesResponse{}, "Success"),
		endpoint.Tags("Images"),
	)

	getImage := endpoint.New("GET", "/images/{image_id}", "Get a specific image",
		endpoint.Description("Get a specific image"),
		endpoint.Handler(controllers.GetImage),
		endpoint.Path("image_id", "string", "uuid", "UUID of an image"),
		endpoint.Response(http.StatusOK, types.Image{}, "Success"),
		endpoint.Tags("Images"),
	)

	createImage := endpoint.New("POST", "/images", "Create an image",
		endpoint.Description("Create an image"),
		endpoint.Consumes("multipart/form-data"),
		endpoint.Handler(controllers.CreateImage),
		endpoint.FormDataMap(map[string]swagger.Parameter{
			"source": {
				Type:        "string",
				Format:      "file",
				Description: "Zip file containing source",
				Required:    true,
			},
			"version": {
				Type:        "string",
				Format:      "string",
				Pattern:     "^(\\d{1,3}\\.?){3}$",
				Description: "Source code semantic version",
				Required:    true,
			},
			"language": {
				Type:        "string",
				Format:      "string",
				Enum:        []string{"Go"},
				Description: "Language the source is written in",
				Required:    true,
			},
			"name": {
				Type:        "string",
				Format:      "string",
				Description: "Name of the service",
				MinLength:   5,
				Required:    true,
			},
		}),
		endpoint.Response(http.StatusNoContent, types.Image{}, "Success"),
		endpoint.Tags("Images"),
	)

	deleteImage := endpoint.New("DELETE", "/images/{image_id}", "Delete an image",
		endpoint.Description("Delete an image"),
		endpoint.Handler(controllers.DeleteImage),
		endpoint.Path("image_id", "string", "uuid", "UUID of an image"),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("Images"),
	)

	return []*swagger.Endpoint{
		getImages,
		getImage,
		createImage,
		deleteImage,
	}
}
