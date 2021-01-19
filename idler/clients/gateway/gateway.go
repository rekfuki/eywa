package gateway

import (
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"

	gwt "eywa/gateway/types"
	"eywa/go-libs/auth"
)

// Client represents gateway client
type Client struct {
	rc *resty.Client
}

// New returns a new gateway client
func New(gatewayURL string) *Client {
	return &Client{
		rc: resty.New().
			SetHostURL(gatewayURL).
			SetLogger(ioutil.Discard).
			SetRetryCount(3).
			SetTimeout(10 * time.Second),
	}
}

// GetFunctions retrieves functions from gateway
func (c *Client) GetFunctions() ([]gwt.FunctionStatusResponse, error) {
	var result gwt.MultiFunctionStatusResponse
	resp, err := c.rc.R().
		SetResult(&result).
		SetHeader("X-User-Id", auth.OperatorUserID).
		SetHeader("X-Real-User-Id", auth.OperatorUserID).
		Get("/eywa/api/system/functions")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		log.Errorf(string(resp.Body()))
		return nil, fmt.Errorf("Gateway responded with unexpected status: %s", resp.Status())
	}

	return result.Objects, nil
}

// ScaleFunction sends scale request to gateway
func (c *Client) ScaleFunction(functionID string, replicas int) error {
	resp, err := c.rc.R().
		SetHeader("X-User-Id", auth.OperatorUserID).
		SetHeader("X-Real-User-Id", auth.OperatorUserID).
		SetPathParams(map[string]string{
			"function_id": functionID,
			"replicas":    fmt.Sprint(replicas),
		}).
		Post("/eywa/api/system/functions/{function_id}/scale/{replicas}")
	if err != nil {
		return err
	}

	if resp.IsError() {
		log.Errorf(string(resp.Body()))
		return fmt.Errorf("Gateway responded with unexpected status: %s", resp.Status())
	}

	return nil
}
