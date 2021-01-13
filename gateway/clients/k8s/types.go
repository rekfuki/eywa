package k8s

import (
	"time"
)

// DeployFunctionRequest represents function deployment creation request
type DeployFunctionRequest struct {
	Image         string
	Service       string
	EnvVars       map[string]string
	Secrets       []Secret
	MinReplicas   int
	MaxReplicas   int
	ScalingFactor int
	Labels        map[string]string
	Annotations   map[string]string
	Limits        *FunctionResources
	Requests      *FunctionResources
}

// FunctionStatus represents the deployed function status in k8s
type FunctionStatus struct {
	Name              string
	Namespace         string
	Image             string
	Env               map[string]string
	MountedSecrets    []string
	Replicas          int
	MaxReplicas       int
	MinReplicas       int
	ScalingFactor     int
	AvailableReplicas int
	Annotations       map[string]string
	Labels            map[string]string
	Limits            *FunctionResources
	Requests          *FunctionResources
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

// FunctionResources represents resources available to the function
type FunctionResources struct {
	Memory string
	CPU    string
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
	Name        string
	Data        map[string]string
	Labels      map[string]string
	Annotations map[string]string
}

// Secret represents k8s secret
type Secret struct {
	Name        string
	MountName   string
	Data        map[string][]byte
	Labels      map[string]string
	Annotations map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
