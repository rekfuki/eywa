package controllers

import (
	"net/http"
	"unsafe"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"eywa/gateway/clients/k8s"
	"eywa/gateway/types"
	"eywa/go-libs/auth"
)

// GetSecrets returns secrets
func GetSecrets(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)

	selector := k8s.LabelSelector().
		Equals(types.UserIDLabel, auth.UserID)

	secrets, err := k8sClient.GetSecretsFiltered(selector)
	if err != nil {
		log.Errorf("Failed to get secrets from k8s: %s", err)
		return err
	}

	secretsResponses := []types.SecretResponse{}
	for _, secret := range secrets {
		secretsResponses = append(secretsResponses, makeSecretResponse(&secret))
	}

	return c.JSON(http.StatusOK, types.MultiSecretResponse{
		Objects: secretsResponses,
		Total:   len(secretsResponses),
	})
}

// GetSecret returns specific secret
func GetSecret(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)
	secretID := c.Param("secret_id")

	selector := k8s.LabelSelector().
		Equals(types.UserIDLabel, auth.UserID).
		Equals(types.SecretIDLabel, secretID)

	secret, err := k8sClient.GetSecretFiltered(selector)
	if err != nil {
		log.Errorf("Failed to get secrets from k8s: %s", err)
		return err
	}

	return c.JSON(http.StatusOK, makeSecretResponse(secret))
}

// CreateSecret creates a new secret inside k8s
func CreateSecret(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)

	var csr types.CreateSecretRequest
	if err := c.Bind(&csr); err != nil {
		return err
	}

	if len(csr.Data) == 0 {
		return c.JSON(http.StatusBadRequest, "Data field cannot be empty")
	}

	if unsafe.Sizeof(csr.Data) > 1024 {
		return c.JSON(http.StatusBadRequest, "Secret size is greater than 1MB")
	}

	secretID := buildK8sName(csr.Name, auth.UserID)
	secretName := secretID[:9] + csr.Name
	selector := k8s.LabelSelector().
		Equals(types.UserIDLabel, auth.UserID).
		Equals(types.SecretIDLabel, secretID)
	secret, err := k8sClient.GetSecretFiltered(selector)
	if err != nil {
		log.Errorf("Failed to get secret from k8s: %s", err)
		return err
	}

	if secret != nil {
		return c.JSON(http.StatusBadRequest, "Secret with specified name already exists")
	}

	sr := &k8s.SecretRequest{
		Name: secretName,
		Data: csr.Data,
		Labels: map[string]string{
			types.UserIDLabel:          auth.UserID,
			types.SecretIDLabel:        secretID,
			types.UserDefinedNameLabel: secretName,
		},
	}

	secret, err = k8sClient.CreateSecret(sr)
	if err != nil {
		log.Errorf("Failed to create secret in k8s: %s", err)
		return err
	}

	return c.JSON(http.StatusCreated, makeSecretResponse(secret))
}

// UpdateSecret updates an existing secret inside k8s
func UpdateSecret(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)
	secretID := c.Param("secret_id")

	var csr types.UpdateSecretRequest
	if err := c.Bind(&csr); err != nil {
		return err
	}

	if len(csr.Data) == 0 {
		return c.JSON(http.StatusBadRequest, "Data field cannot be empty")
	}

	if unsafe.Sizeof(csr.Data) > 1024 {
		return c.JSON(http.StatusBadRequest, "Secret size is greater than 1MB")
	}

	selector := k8s.LabelSelector().
		Equals(types.UserIDLabel, auth.UserID).
		Equals(types.SecretIDLabel, secretID)
	secret, err := k8sClient.GetSecretFiltered(selector)
	if err != nil {
		log.Errorf("Failed to get secret from k8s: %s", err)
		return err
	}

	if secret == nil {
		return c.JSON(http.StatusNotFound, "Secret Not Found")
	}

	sr := &k8s.SecretRequest{
		Data: csr.Data,
	}

	secret, err = k8sClient.UpdateSecret(secretID, sr)
	if err != nil {
		log.Errorf("Failed to create secret in k8s: %s", err)
		return err
	}

	return c.JSON(http.StatusCreated, makeSecretResponse(secret))
}

// DeleteSecret deletes a secret from k8s
func DeleteSecret(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	k8sClient := c.Get("k8s").(*k8s.Client)
	secretID := c.Param("secret_id")

	selector := k8s.LabelSelector().
		Equals(types.UserIDLabel, auth.UserID).
		Equals(types.SecretIDLabel, secretID)
	secret, err := k8sClient.GetSecretFiltered(selector)
	if err != nil {
		log.Errorf("Failed to get secret from k8s: %s", err)
		return err
	}

	if secret == nil {
		return c.JSON(http.StatusNotFound, "Secret Not Found")
	}

	if err := k8sClient.DeleteSecret(secretID); err != nil {
		log.Errorf("Failed to delete a secret")
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func makeSecretResponse(s *k8s.Secret) types.SecretResponse {
	secret := types.SecretResponse{
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}

	if val, exists := s.Labels[types.SecretIDLabel]; exists {
		secret.ID = val
	}

	return secret
}
