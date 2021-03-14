package types

import (
	"encoding/json"
)

// QueryRequest represents prometheus query request
type QueryRequest struct {
	Type          string   `json:"type" enum:"instant,range" binding:"required"`
	Series        []string `json:"series" binding:"required"`
	LabelMatchers string   `json:"label_matchers"`
	GroupBy       string   `json:"group_by"`
	Query         string   `json:"query" binding:"required"`
	Start         float64  `json:"start"`
	End           float64  `json:"end"`
	Step          int64    `json:"step"`
}

// QueryResponse represents prometheus  query response
type QueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string          `json:"resultType"`
		Result     json.RawMessage `json:"result"`
	}
}
