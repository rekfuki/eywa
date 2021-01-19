package types

// SystemScaleFunctionRequest represents function scale request payload
type SystemScaleFunctionRequest struct {
	Name     string `json:"name" binding:"required"`
	Replicas int    `json:"replicas" binding:"required"`
}
