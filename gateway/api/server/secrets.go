package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"eywa/gateway/api/controllers"
	"eywa/gateway/types"
)

func secretsAPI() []*swagger.Endpoint {
	getSecrets := endpoint.New("GET", "/secrets", "Get secrets",
		endpoint.Description("Get all secrets belonging to a user"),
		endpoint.Handler(controllers.GetSecrets),
		endpoint.Response(http.StatusOK, types.MultiSecretResponse{}, "Success"),
		endpoint.Tags("Secrets"),
	)

	getSecret := endpoint.New("GET", "/secrets/{secret_id}", "Get specific secret",
		endpoint.Description("Get a secret belonging to a user"),
		endpoint.Handler(controllers.GetSecret),
		endpoint.Path("secret_id", "string", "uuid", "UUID of a secret"),
		endpoint.Response(http.StatusOK, types.SecretResponse{}, "Success"),
		endpoint.Tags("Secrets"),
	)

	createSecret := endpoint.New("POST", "/secrets", "Create a secret",
		endpoint.Description("Create a function from an image"),
		endpoint.Handler(controllers.CreateSecret),
		endpoint.Body(types.CreateSecretRequest{}, "Secret creation payload", true),
		endpoint.Response(http.StatusCreated, types.SecretResponse{}, "Success"),
		endpoint.Tags("Secrets"),
	)

	updateSecret := endpoint.New("PUT", "/secrets/{secret_id}", "Update a secret",
		endpoint.Description("Update a secret. All data is overriden"),
		endpoint.Handler(controllers.UpdateSecret),
		endpoint.Path("secret_id", "string", "uuid", "UUID of a secret"),
		endpoint.Body(types.UpdateSecretRequest{}, "Secret update payload", true),
		endpoint.Response(http.StatusOK, types.SecretResponse{}, "Success"),
		endpoint.Tags("Secrets"),
	)

	deleteSecret := endpoint.New("DELETE", "/secrets/{secret_id}", "Delete a secret",
		endpoint.Description("Delete a secret from k8s"),
		endpoint.Handler(controllers.DeleteSecret),
		endpoint.Path("secret_id", "string", "uuid", "UUID of a secret"),
		endpoint.Response(http.StatusNoContent, "", "Success"),
		endpoint.Tags("Secrets"),
	)

	return []*swagger.Endpoint{
		getSecrets,
		getSecret,
		createSecret,
		updateSecret,
		deleteSecret,
	}
}
