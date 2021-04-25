package prometheus

// VectorQueryResponse represents prometheus vector query response
type VectorQueryResponse struct {
	Data struct {
		Result []struct {
			Metric struct {
				Code       string `json:"code"`
				FunctionID string `json:"function_id"`
			}
			Value []interface{} `json:"value"`
		}
	}
}
