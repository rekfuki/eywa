package types

// CreateFunctionRequest represents function deployment creation request
type CreateFunctionRequest struct {
	Image       string             `json:"image" binding:"required"`
	Service     string             `json:"service" binding:"required"`
	EnvVars     map[string]string  `json:"env_vars"`
	Secrets     []string           `json:"secrets"`
	Labels      *map[string]string `json:"labels"`
	Annotations *map[string]string `json:"annotations"`
	Limits      *FunctionResources `json:"limits"`
	Requests    *FunctionResources `json:"requests"`
}

// FunctionResources represents resources available to the function
type FunctionResources struct {
	Memory string `json:"memory"`
	CPU    string `json:"cpu"`
}

// DeleteFunctionRequest represents function deployment deletion payload
type DeleteFunctionRequest struct {
	Name string `json:"name" binding:"required"`
}

// ScaleFunctionRequest represents function scale request payload
type ScaleFunctionRequest struct {
	Name     string `json:"name" binding:"required"`
	Replicas int32  `json:"replicas" binding:"required"`
}
