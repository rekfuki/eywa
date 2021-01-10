package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Secret represents k8s secret
type Secret struct {
	Name string
	Data map[string][]byte
}

// GetSecrets returns specified secrets if they exist
func (c *Client) GetSecrets(secrets []string) ([]Secret, error) {
	kube := c.clientset.CoreV1().Secrets(faasNamespace)
	opts := metav1.GetOptions{}

	s := []Secret{}
	for _, sn := range secrets {
		secret, err := kube.Get(context.TODO(), sn, opts)
		if err != nil {
			return nil, err
		}

		if secret == nil {
			return s, fmt.Errorf("Secret %q not found in k8s", sn)
		}
		s = append(s, Secret{
			Name: secret.Name,
			Data: secret.Data,
		})
	}

	return s, nil
}
