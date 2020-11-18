package k8s

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"

	"eywa/gateway/types"
)

type Config struct {
}

type Client struct {
	clientset *kubernetes.Clientset
}

// Setup initialises a k8s client
func Setup(conf *Config) (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{clientset}, nil
}

type Deployment struct {
}

func (c *Client) CreateDeployment(fn *types.Function) (*Deployment, error) {
	deploymentsClient := c.clientset.AppsV1().Deployments("faas")

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo-deployment",
			Namespace: fn.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"function": fn.Name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: fn.Namespace,
					Labels: map[string]string{
						"function": fn.Name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  fn.Name,
							Image: fn.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}

	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	spew.Dump(result)

	return &Deployment{}, nil
}
