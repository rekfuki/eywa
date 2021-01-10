package types

import (
	"time"
)

// DeployFunctionRequest represents a request payload for function deployment
type DeployFunctionRequest struct {
	ImageID       string            `json:"image_id" format:"uuid" binding:"required"`
	Name          string            `json:"name" min_length:"5" pattern:"^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$" binding:"required"`
	EnvVars       map[string]string `json:"env_vars" min_length:"1"`
	Secrets       []string          `json:"secrets"`
	MinReplicas   int               `json:"min_replicas" minimum:"0" maximum:"100" binding:"required"`
	MaxReplicas   int               `json:"max_replicas" minimum:"1" maximum:"100" binding:"required"`
	ScalingFactor int               `json:"scaling_factor" minimum:"0" maximum:"100"`
	MaxInflight   int               `json:"max_concurrency" minimum:"0"`
	WriteDebug    bool              `json:"write_debug"`
	ReadTimeout   time.Duration     `json:"read_timeout" pattern:"^[1-9]{1}\\d{0,}s$"`
	WriteTimeout  time.Duration     `json:"write_timeout" pattern:"^[1-9]{1}\\d{0,}s$"`
	Resources     FunctionResources `json:"resources" binding:"required"`
}

// FunctionStatusResponse represents a function status response that has been sanitized
type FunctionStatusResponse struct {
	ID                string            `json:"id"`
	ImageID           string            `json:"image_id"`
	ShortName         string            `json:"short_name"`
	FullName          string            `json:"full_name"`
	EnvVars           map[string]string `json:"env_vars"`
	MountedSecrets    []string          `json:"mounted_secrets"`
	AvailableReplicas int32             `json:"available_replicas"`
	MinReplicas       int32             `json:"min_replicas"`
	MaxReplicas       int32             `json:"max_replicas"`
	ScalingFactor     int32             `json:"scaling_factor" minimum:"0" maximum:"100"`
	MaxInflight       int               `json:"max_concurrency" minimum:"0"`
	WriteDebug        bool              `json:"write_debug"`
	ReadTimeout       time.Duration     `json:"read_timeout"`
	WriteTimeout      time.Duration     `json:"write_timeout"`
	Resources         FunctionResources `json:"resources"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// FunctionResources represents request and limit resources of k8s
type FunctionResources struct {
	MinCPU    string `json:"min_cpu" pattern:"^[1-9]{1}\\d{0,}m$" binding:"required"`
	MinMemory string `json:"min_memory" pattern:"^[1-9]{1}\\d{0,}Mi$" binding:"required"`
	MaxCPU    string `json:"max_cpu" pattern:"^[1-9]{1}\\d{0,}m$" binding:"required"`
	MaxMemory string `json:"max_memory" pattern:"^[1-9]{1}\\d{0,}Mi$" binding:"required"`
}
