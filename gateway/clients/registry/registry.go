package registry

import (
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"

	"eywa/go-libs/auth"
	rt "eywa/registry/types"
)

// Client represents registry client
type Client struct {
	rc *resty.Client
}

// New returns a new registry client
func New(registryURL string) *Client {
	return &Client{
		rc: resty.New().
			SetHostURL(registryURL).
			SetLogger(ioutil.Discard).
			SetRetryCount(3).
			SetTimeout(10 * time.Second),
	}
}

// GetImage retrieves image from registry filtered by user and the id
func (c *Client) GetImage(imageID, userID string) (*rt.Image, error) {
	var result rt.Image
	resp, err := c.rc.R().
		SetResult(&result).
		SetHeader("X-Eywa-User-Id", userID).
		SetHeader("X-Eywa-Real-User-Id", auth.OperatorUserID).
		Get("/eywa/api/images/" + imageID)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		log.Errorf(string(resp.Body()))
		return nil, fmt.Errorf("Registry responded with unexpected status: %s", resp.Status())
	}

	return &result, nil
}
