package prometheus

import (
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

// Client represents prometheus client
type Client struct {
	rc *resty.Client
}

// New returns a new prometheus client
func New(prometheusURL string) *Client {
	return &Client{
		rc: resty.New().
			SetHostURL(prometheusURL + "/api/v1").
			SetLogger(ioutil.Discard).
			SetRetryCount(3).
			SetTimeout(10 * time.Second),
	}
}

// QueryMetrics queries prometheus for metrics
func (c *Client) QueryMetrics(query string) (*VectorQueryResponse, error) {
	var result VectorQueryResponse
	resp, err := c.rc.R().
		SetResult(&result).
		SetQueryParam("query", query).
		Get("/query")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		log.Errorf(string(resp.Body()))
		return nil, fmt.Errorf("Prometheus responded with unexpected status: %s", resp.Status())
	}

	return &result, nil
}
