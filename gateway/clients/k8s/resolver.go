package k8s

import (
	"fmt"
	"math/rand"
	"strings"
)

func (c *Client) Resolve(fnName string) (string, error) {
	if strings.Contains(fnName, ".") {
		fnName = strings.TrimSuffix(fnName, "."+faasNamespace)
	}

	svc, err := c.endpointLister.Get(fnName)
	if err != nil {
		return "", fmt.Errorf("Error listing \"%s.%s\": %s", fnName, faasNamespace, err)
	}

	if len(svc.Subsets) == 0 {
		return "", fmt.Errorf("No subsets available for \"%s.%s\"", fnName, faasNamespace)
	}

	all := len(svc.Subsets[0].Addresses)
	if len(svc.Subsets[0].Addresses) == 0 {
		return "", fmt.Errorf("No addresses in subset for \"%s.%s\"", fnName, faasNamespace)
	}

	target := rand.Intn(all)

	serviceIP := svc.Subsets[0].Addresses[target].IP

	return fmt.Sprintf("http://%s:%d", serviceIP, 8080), nil
}
