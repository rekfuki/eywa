package metrics

import (
	"encoding/json"
	"eywa/gateway/types"
)

// QueryMetrics queries prometheus for metrics
func (c *Client) QueryMetrics(endpoint, query string) (*types.QueryResponse, *types.Error) {
	var result types.QueryResponse
	resp, err := c.promrc.R().
		SetResult(&result).
		SetQueryString(query).
		Get(endpoint)
	if err != nil {
		return nil, types.SystemError(err.Error())
	}

	if resp.IsError() {
		var promErr struct {
			Error string `json:"error"`
		}

		if err := json.Unmarshal(resp.Body(), &promErr); err != nil {
			return nil, types.SystemError(err.Error())
		}

		if resp.StatusCode() == 400 {
			return nil, types.BadRequestError(promErr.Error)
		}
		return nil, types.SystemError(promErr.Error)
	}

	return &result, nil
}
