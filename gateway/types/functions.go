package types

import (
	"time"
)

// FunctionRequest reprsents shared fields between update and deploy
type FunctionRequest struct {
	EnvVars       map[string]string `json:"env_vars" min_length:"1"`
	Secrets       []string          `json:"secrets" unique_items:"true" format:"uuid"`
	MinReplicas   int               `json:"min_replicas" minimum:"0" maximum:"100" binding:"required"`
	MaxReplicas   int               `json:"max_replicas" minimum:"1" maximum:"100" binding:"required"`
	ScalingFactor int               `json:"scaling_factor" minimum:"0" maximum:"100"`
	MaxInflight   int               `json:"max_concurrency" minimum:"0"`
	WriteDebug    bool              `json:"write_debug"`
	ReadTimeout   string            `json:"read_timeout" pattern:"^[1-9]{1}\\d{0,}s$"`
	WriteTimeout  string            `json:"write_timeout" pattern:"^[1-9]{1}\\d{0,}s$"`
	Resources     FunctionResources `json:"resources" binding:"required"`
}

// DeployFunctionRequest represents a request payload for function deployment
type DeployFunctionRequest struct {
	ImageID string `json:"image_id" format:"uuid" binding:"required"`
	Name    string `json:"name" min_length:"5" pattern:"^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$" binding:"required"`
	FunctionRequest
}

// UpdateFunctionRequest represents a request payload for function deployment update
type UpdateFunctionRequest struct {
	FunctionRequest
}

// MultiFunctionStatusResponse represents the response of multiple status returns
type MultiFunctionStatusResponse struct {
	Objects []FunctionStatusResponse `json:"objects"`
	Total   int                      `json:"total"`
}

// FunctionStatusResponse represents a function status response that has been sanitized
type FunctionStatusResponse struct {
	ID                string            `json:"id"`
	ImageID           string            `json:"image_id"`
	Name              string            `json:"short_name"`
	EnvVars           map[string]string `json:"env_vars" min_length:"1"`
	Secrets           []string          `json:"secrets"`
	AvailableReplicas int               `json:"available_replicas"`
	MinReplicas       int               `json:"min_replicas" minimum:"0" maximum:"100" binding:"required"`
	MaxReplicas       int               `json:"max_replicas" minimum:"1" maximum:"100" binding:"required"`
	ScalingFactor     int               `json:"scaling_factor" minimum:"0" maximum:"100"`
	MaxInflight       int               `json:"max_concurrency" minimum:"0"`
	WriteDebug        bool              `json:"write_debug"`
	ReadTimeout       string            `json:"read_timeout" pattern:"^[1-9]{1}\\d{0,}s$"`
	WriteTimeout      string            `json:"write_timeout" pattern:"^[1-9]{1}\\d{0,}s$"`
	Resources         FunctionResources `json:"resources" binding:"required"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	DeletedAt         *time.Time        `json:"deleted_at,omitempty"`
}

// FunctionResources represents request and limit resources of k8s
type FunctionResources struct {
	MinCPU    string `json:"min_cpu" pattern:"^[1-9]{1}\\d{0,}m$" binding:"required"`
	MinMemory string `json:"min_memory" pattern:"^[1-9]{1}\\d{0,}Mi$" binding:"required"`
	MaxCPU    string `json:"max_cpu" pattern:"^[1-9]{1}\\d{0,}m$" binding:"required"`
	MaxMemory string `json:"max_memory" pattern:"^[1-9]{1}\\d{0,}Mi$" binding:"required"`
}
