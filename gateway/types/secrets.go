package types

import "time"

// CreateSecretRequest represents a request to create a secret
type CreateSecretRequest struct {
	Name string            `json:"name" min_length:"5" pattern:"^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$" binding:"required"`
	Data map[string]string `json:"data" binding:"required"`
}

// UpdateSecretRequest represents a request to update a secret
type UpdateSecretRequest struct {
	Data map[string]string `json:"data" binding:"required"`
}

// MultiSecretResponse represents the response of multiple secrets
type MultiSecretResponse struct {
	Objects []SecretResponse `json:"objects"`
	Total   int              `json:"total"`
}

// SecretResponse represents the responses involving secrets
type SecretResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
