package gateway

import (
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"

	gt "eywa/gateway/types"
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

// CreateDatabaseSecret calls gateway to create database secret for user to user
func (c *Client) CreateDatabaseSecret(userID string, request gt.CreateSecretRequest) error {
	resp, err := c.rc.R().
		SetBody(request).
		SetHeader("X-Eywa-User-Id", userID).
		SetHeader("X-Eywa-Real-User-Id", auth.OperatorUserID).
		Post("/eywa/api/secrets")
	if err != nil {
		return err
	}

	if resp.IsError() {
		log.Errorf(string(resp.Body()))
		return fmt.Errorf("Gateway responded with unexpected status: %s", resp.Status())
	}

	return nil
}
