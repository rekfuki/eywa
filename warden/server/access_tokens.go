package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"eywa/warden/controllers"
	"eywa/warden/types"
)

func accessTokensAPI() []*swagger.Endpoint {
	getAccessTokens := endpoint.New("GET", "/tokens", "Get access tokens",
		endpoint.Description("Get all user access tokens"),
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
			"query": {
				Type:        "string",
				Description: "Query string to search by",
			},
		}),
		endpoint.Handler(controllers.GetAccessTokens),
		endpoint.Response(http.StatusOK, types.AccessTokensResponse{}, "Success"),
		endpoint.Tags("Access Tokens"),
	)

	createAccessToken := endpoint.New("POST", "/tokens", "Create an access token",
		endpoint.Description("Create an access token to be used with the api"),
		endpoint.Handler(controllers.CreateToken),
		endpoint.Body(types.CreateTokenRequest{}, "Token creation payload", true),
		endpoint.Response(http.StatusCreated, types.AccessToken{}, "Success"),
		endpoint.Tags("Access Tokens"),
	)

	deleteAccessToken := endpoint.New("DELETE", "/tokens/{token_id}", "Delete access token",
		endpoint.Description("Delete user access token"),
		endpoint.Handler(controllers.DeleteAccessToken),
		endpoint.Path("token_id", "string", "uuid", "UUID of a token"),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("Access Tokens"),
	)

	return []*swagger.Endpoint{
		createAccessToken,
		getAccessTokens,
		deleteAccessToken,
	}
}
