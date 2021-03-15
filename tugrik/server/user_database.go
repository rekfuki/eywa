package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"eywa/tugrik/controllers"
	"eywa/tugrik/types"
)

func userDatabaseAPI() []*swagger.Endpoint {
	getUserDatabase := endpoint.New("GET", "/database", "Get user database",
		endpoint.Description("Create user database"),
		endpoint.Handler(controllers.GetUserDatabase),
		endpoint.Response(http.StatusOK, types.UserDatabaseInfo{}, "Success"),
		endpoint.Tags("User Database"),
	)

	createUserDatabase := endpoint.New("POST", "/system/database/{user_id}", "Create user database",
		endpoint.Description("Create user database"),
		endpoint.Path("user_id", "string", "uuid", "User UUID"),
		endpoint.Handler(controllers.CreateUserDatabase),
		endpoint.Response(http.StatusCreated, types.CreateUserDatabaseResponse{}, "Success"),
		endpoint.Tags("User Database"),
	)

	deleteCollection := endpoint.New("DELETE", "/database/collections/{collection_name}", "Delete a collection",
		endpoint.Description("Delete a collection from user database"),
		endpoint.Path("collection_name", "string", "string", "Collection name"),
		endpoint.Handler(controllers.DeleteCollection),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("User Database"),
	)

	return []*swagger.Endpoint{
		createUserDatabase,
		getUserDatabase,
		deleteCollection,
	}
}
