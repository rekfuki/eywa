package k8s

import (
	"time"
)

// DeployFunctionRequest represents function deployment creation request
type DeployFunctionRequest struct {
	Image         string             `json:"image" binding:"required"`
	Service       string             `json:"service" binding:"required"`
	EnvVars       map[string]string  `json:"env_vars"`
	Secrets       []string           `json:"secrets"`
	MinReplicas   int                `json:"min_replicas" minimum:"1" maximum:"100" default:"1"`
	MaxReplicas   int                `json:"max_replicas" minimum:"1" maximum:"100" default:"100"`
	ScalingFactor int                `json:"scaling_factor" minimum:"0" maximum:"100" default:"20"`
	Labels        map[string]string  `json:"labels"`
	Annotations   map[string]string  `json:"annotations"`
	Limits        *FunctionResources `json:"limits"`
	Requests      *FunctionResources `json:"requests"`
}

// FunctionStatus represents the deployed function status in k8s
type FunctionStatus struct {
	Name              string             `json:"name"`
	Namespace         string             `json:"namespace"`
	Image             string             `json:"image"`
	Env               map[string]string  `json:"env"`
	MountedSecrets    []string           `json:"mounted_secrets"`
	Replicas          int                `json:"replicas"`
	MaxReplicas       int                `json:"max_replicas"`
	MinReplicas       int                `json:"min_replicas"`
	ScalingFactor     int                `json:"scaling_factor"`
	AvailableReplicas int                `json:"available_replicas"`
	Annotations       map[string]string  `json:"annotations"`
	Labels            map[string]string  `json:"labels"`
	Limits            *FunctionResources `json:"limits"`
	Requests          *FunctionResources `json:"requests"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	DeletedAt         *time.Time         `json:"deleted_at,omitempty"`
}

// FunctionResources represents resources available to the function
type FunctionResources struct {
	Memory string `json:"memory"`
	CPU    string `json:"cpu"`
}

// FunctionZeroScaleResult represents status of function scaling from zero
type FunctionZeroScaleResult struct {
	Found     bool
	Available bool
	Duration  time.Duration
}

// ResourceLimits represents response of resource limits
type ResourceLimits struct {
	MinCPU string
	MaxCPU string
	MinMem string
	MaxMem string
}

// SecretRequest represents a secret creation/update request
type SecretRequest struct {
	Name        string            `json:"name"`
	Data        map[string]string `json:"data"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// Secret represents k8s secret
type Secret struct {
	Name        string
	Data        map[string][]byte
	Labels      map[string]string
	Annotations map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
