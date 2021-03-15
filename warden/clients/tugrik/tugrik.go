package tugrik

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"

	"eywa/go-libs/auth"
)

// Client represents tugrik client
type Client struct {
	rc *resty.Client
}

// New returns a new tugrik client
func New(tugrikURL string) *Client {
	return &Client{
		rc: resty.New().
			SetHostURL(tugrikURL).
			SetLogger(ioutil.Discard).
			SetRetryCount(3).
			SetTimeout(10 * time.Second),
	}
}

// CreateUserDatabase calls Tugrik to create a database for the user
// We don't care about the response as long as it's not an error including 409
func (c *Client) CreateUserDatabase(userID string) error {
	resp, err := c.rc.R().
		SetHeader("X-Eywa-User-Id", auth.OperatorUserID).
		SetHeader("X-Eywa-Real-User-Id", auth.OperatorUserID).
		Post("/eywa/api/system/database/" + userID)
	if err != nil {
		return err
	}

	if resp.IsError() {
		if resp.StatusCode() == http.StatusConflict {
			return nil
		}
		log.Errorf(string(resp.Body()))
		return fmt.Errorf("Tugrik responded with unexpected status: %s", resp.Status())
	}

	return nil
}
