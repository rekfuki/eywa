package k8s

import (
	"context"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

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

// GetSecretsFiltered returns specified secrets filtered by labels
func (c *Client) GetSecretsFiltered(filter Selector) ([]Secret, error) {
	return c.getSecrets(filter...)
}

// GetSecretFiltered returns a secret that is filtered by labels
func (c *Client) GetSecretFiltered(filter Selector) (*Secret, error) {
	return c.getSecret(filter...)
}

func (c *Client) getSecret(r ...labels.Requirement) (*Secret, error) {
	secrets, err := c.listSecrets(r...)
	if err != nil {
		return nil, err
	}

	if len(secrets.Items) == 0 {
		return nil, nil
	}

	if len(secrets.Items) != 1 {
		log.Warnf("K8s returned more than one result when only one was expected: %#v", r)
	}

	return convertSecret(&secrets.Items[0])
}

func (c *Client) getSecrets(r ...labels.Requirement) ([]Secret, error) {
	secrets, err := c.listSecrets(r...)
	if err != nil {
		return nil, err
	}

	convertedSecrets := []Secret{}
	for _, secret := range secrets.Items {
		cs, err := convertSecret(&secret)
		if err != nil {
			return nil, err
		}
		convertedSecrets = append(convertedSecrets, *cs)
	}

	return convertedSecrets, err
}

func (c *Client) listSecrets(r ...labels.Requirement) (*corev1.SecretList, error) {
	opts := metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind: "Secret",
		},
		LabelSelector: labels.NewSelector().Add(r...).String(),
	}

	return c.clientset.CoreV1().Secrets(faasNamespace).List(context.TODO(), opts)
}

// CreateSecret creates a new secret inside k8s
func (c *Client) CreateSecret(sr *SecretRequest) (*Secret, error) {
	data := map[string][]byte{}

	if len(sr.Data) == 0 {
		return nil, fmt.Errorf("Secret data must not be empty")
	}

	for k, v := range sr.Data {
		data[k] = []byte(v)
	}

	sr.Labels[updatedAtLabel] = fmt.Sprint(time.Now().Unix())

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        sr.Name,
			Namespace:   faasNamespace,
			Annotations: sr.Annotations,
			Labels:      sr.Labels,
		},
		Type: v1.SecretTypeOpaque,
		Data: data,
	}

	secret, err := c.clientset.CoreV1().
		Secrets(faasNamespace).
		Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return convertSecret(secret)
}

// UpdateSecret updates an existing secret inside k8s
// Only data field has to be set, everything else is ignored in the request
func (c *Client) UpdateSecret(secretID string, sr *SecretRequest) (*Secret, error) {
	secret, err := c.clientset.CoreV1().
		Secrets(faasNamespace).
		Get(context.TODO(), secretID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	data := map[string][]byte{}

	if len(sr.Data) == 0 {
		return nil, fmt.Errorf("Secret data must not be empty")
	}

	for k, v := range sr.Data {
		data[k] = []byte(v)
	}

	secret.Data = data
	secret.Labels[updatedAtLabel] = fmt.Sprint(time.Now().Unix())

	secret, err = c.clientset.CoreV1().
		Secrets(faasNamespace).
		Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return convertSecret(secret)
}

// DeleteSecret deletes a secret from k8s
// If secret is not found, no error is returned
func (c *Client) DeleteSecret(secretID string) error {
	err := c.clientset.CoreV1().
		Secrets(faasNamespace).
		Delete(context.TODO(), secretID, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func convertSecret(secret *v1.Secret) (*Secret, error) {
	s := &Secret{
		Name:        secret.Name,
		Data:        secret.Data,
		Labels:      secret.Labels,
		Annotations: secret.Annotations,
		CreatedAt:   secret.CreationTimestamp.Time,
	}

	for k, v := range secret.Labels {
		if k == updatedAtLabel {
			i64, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			s.UpdatedAt = time.Unix(i64, 0)
		}
	}

	return s, nil
}
