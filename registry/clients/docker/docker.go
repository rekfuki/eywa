package docker

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/resty.v1"
)

// Client represents docker client
type Client struct {
	rc *resty.Client
}

// New creates a new client
func New(registry, user, pass string) *Client {
	return &Client{
		rc: resty.New().
			SetHostURL(fmt.Sprintf("https://%s/v2", registry)).
			SetLogger(ioutil.Discard).
			SetRetryCount(3).
			SetTimeout(10*time.Second).
			SetBasicAuth(user, pass).
			SetHeader("Accept", "application/vnd.docker.distribution.manifest.v2+json"),
	}
}

// DeleteImage deletes docker image from the registry
func (c *Client) DeleteImage(repository, version string) error {
	resp, err := c.rc.R().Get(fmt.Sprintf("/%s/manifests/%s", repository, version))
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("Docker registry returned unxpected response: %s", resp.Status())
	}

	digest := resp.Header().Get("Docker-Content-Digest")
	if digest == "" {
		return fmt.Errorf("Missing Docker-Content-Digest header")
	}

	resp, err = c.rc.R().Delete(fmt.Sprintf("/%s/manifests/%s", repository, digest))
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("Docker registry returned unxpected response: %s", resp.Status())
	}

	return nil
}
