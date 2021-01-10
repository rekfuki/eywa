package types

// SystemDeleteFunctionRequest represents function deployment deletion payload
type SystemDeleteFunctionRequest struct {
	Name string `json:"name" binding:"required"`
}

// SystemScaleFunctionRequest represents function scale request payload
type SystemScaleFunctionRequest struct {
	Name     string `json:"name" binding:"required"`
	Replicas int32  `json:"replicas" binding:"required"`
}

// // GetFunctionsFilteredRequest represents the payload used to query k8s for functions
// type GetFunctionsFilteredRequest struct {
// 	Labels map[string]string `json:"labels" binding:"required"`
// }

// // GetFunctionsFilteredResponse represents the response of a query
// type GetFunctionsFilteredResponse struct {
// 	Objects []FunctionStatus `json:"objects"`
// 	Total   int              `json:"total"`
// }
