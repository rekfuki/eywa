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
		secretsResponses = append(secretsResponses, makeSecretResponse(&secret, nil))
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

	// TODO: This might have a substantial performance hit
	// since we have to iterate through every single user deployment
	// and their mounted secrets.
	// investigate and figure out a better solution
	filter := k8s.LabelSelector().Equals(types.UserIDLabel, auth.UserID)
	fss, err := k8sClient.GetFunctionsStatusFiltered(filter)
	if err != nil {
		log.Errorf("Failed to get functions from k8s: ", err)
		return err
	}

	mounts := []types.MountedFunction{}
	for _, fs := range fss {
		for _, ms := range fs.MountedSecrets {
			if ms == secret.Name {
				name, exists := fs.Labels[types.UserDefinedNameLabel]
				if !exists {
					continue
				}
				mounts = append(mounts, types.MountedFunction{
					ID:   fs.Name,
					Name: name,
				})
			}
		}
	}

	return c.JSON(http.StatusOK, makeSecretResponse(secret, mounts))
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
			types.UserDefinedNameLabel: csr.Name,
			types.SecretNameLabel:      secretName,
		},
	}

	secret, err = k8sClient.CreateSecret(sr)
	if err != nil {
		log.Errorf("Failed to create secret in k8s: %s", err)
		return err
	}

	return c.JSON(http.StatusCreated, makeSecretResponse(secret, nil))
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

	if len(csr.Upserts) == 0 && len(csr.Deletes) == 0 {
		return c.JSON(http.StatusBadRequest, "Either deletes or updates must be provided")
	}

	if unsafe.Sizeof(csr.Upserts) > 1024 && unsafe.Sizeof(csr.Deletes) > 1024 {
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

	for _, toDelete := range csr.Deletes {
		delete(secret.Data, toDelete)
	}

	for k, data := range csr.Upserts {
		secret.Data[k] = []byte(data)
	}

	secret, err = k8sClient.UpdateSecret(secret.Name, secret.Data)
	if err != nil {
		log.Errorf("Failed to update secret in k8s: %s", err)
		return err
	}

	return c.JSON(http.StatusCreated, makeSecretResponse(secret, nil))
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

	if err := k8sClient.DeleteSecret(secret.Name); err != nil {
		log.Errorf("Failed to delete a secret")
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func makeSecretResponse(s *k8s.Secret, mounts []types.MountedFunction) types.SecretResponse {
	secret := types.SecretResponse{
		Name:       s.Name,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
		Mounts:     mounts,
		DataFields: []string{},
	}

	for k := range s.Data {
		secret.DataFields = append(secret.DataFields, k)
	}

	if val, exists := s.Labels[types.SecretIDLabel]; exists {
		secret.ID = val
	}

	return secret
}
