package k8s

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"eywa/gateway/types"
)

var replicaCount = int32(1)

func (c *Client) CreateFunction(request *types.CreateFunctionRequest, secrets []Secret) error {
	envVars := []corev1.EnvVar{}
	for k, v := range request.EnvVars {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	labels := map[string]string{
		faasIDLabel: request.Service,
	}

	if request.Labels != nil {
		for k, v := range *request.Labels {
			labels[k] = v
		}
	}

	resources := &apiv1.ResourceRequirements{
		Limits:   apiv1.ResourceList{},
		Requests: apiv1.ResourceList{},
	}

	// Set Memory limits
	if request.Limits != nil && len(request.Limits.Memory) > 0 {
		qty, err := resource.ParseQuantity(request.Limits.Memory)
		if err != nil {
			return err
		}
		resources.Limits[apiv1.ResourceMemory] = qty
	}

	if request.Requests != nil && len(request.Requests.Memory) > 0 {
		qty, err := resource.ParseQuantity(request.Requests.Memory)
		if err != nil {
			return err
		}
		resources.Requests[apiv1.ResourceMemory] = qty
	}

	// Set CPU limits
	if request.Limits != nil && len(request.Limits.CPU) > 0 {
		qty, err := resource.ParseQuantity(request.Limits.CPU)
		if err != nil {
			return err
		}
		resources.Limits[apiv1.ResourceCPU] = qty
	}

	if request.Requests != nil && len(request.Requests.CPU) > 0 {
		qty, err := resource.ParseQuantity(request.Requests.CPU)
		if err != nil {
			return err
		}
		resources.Requests[apiv1.ResourceCPU] = qty
	}

	imagePullPolicy := apiv1.PullIfNotPresent

	annotations := map[string]string{}
	if request.Annotations != nil {
		annotations = *request.Annotations
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        request.Service,
			Annotations: annotations,
			Labels: map[string]string{
				faasIDLabel: request.Service,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					faasIDLabel: request.Service,
				},
			},
			Replicas: &replicaCount,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(0),
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(1),
					},
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        request.Service,
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  request.Service,
							Image: request.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env:             envVars,
							Resources:       *resources,
							ImagePullPolicy: imagePullPolicy,
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
					DNSPolicy:     corev1.DNSClusterFirst,
				},
			},
		},
	}

	secretVolumeProjections := []apiv1.VolumeProjection{}
	for _, secret := range secrets {
		projectedPaths := []apiv1.KeyToPath{}
		for secretKey := range secret.Data {
			projectedPaths = append(projectedPaths, apiv1.KeyToPath{Key: secretKey, Path: secretKey})
		}

		projection := &apiv1.SecretProjection{Items: projectedPaths}
		projection.Name = secret.Name
		secretProjection := apiv1.VolumeProjection{
			Secret: projection,
		}
		secretVolumeProjections = append(secretVolumeProjections, secretProjection)
	}

	if len(secretVolumeProjections) > 0 {
		volumeName := fmt.Sprintf("%s-projected-secrets", request.Service)
		projectedSecrets := apiv1.Volume{
			Name: volumeName,
			VolumeSource: apiv1.VolumeSource{
				Projected: &apiv1.ProjectedVolumeSource{
					Sources: secretVolumeProjections,
				},
			},
		}

		deployment.Spec.Template.Spec.Volumes = []corev1.Volume{projectedSecrets}

		for i := range deployment.Spec.Template.Spec.Containers {
			mount := apiv1.VolumeMount{
				Name:      volumeName,
				ReadOnly:  true,
				MountPath: faasSecretMount,
			}

			deployment.Spec.Template.Spec.Containers[i].VolumeMounts = []apiv1.VolumeMount{mount}
		}
	}

	_, err := c.clientset.AppsV1().
		Deployments(faasNamespace).
		Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        request.Service,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				faasIDLabel: request.Service,
			},
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: corev1.ProtocolTCP,
					Port:     8080,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8080,
					},
				},
			},
		},
	}
	_, err = c.clientset.CoreV1().
		Services(faasNamespace).
		Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteFunction(fn *Function) error {
	foregroundPolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &foregroundPolicy}

	if err := c.clientset.AppsV1().
		Deployments(fn.Namespace).
		Delete(context.TODO(), fn.Name, *opts); err != nil {
		return err
	}

	if err := c.clientset.CoreV1().
		Services(fn.Namespace).
		Delete(context.TODO(), fn.Name, *opts); err != nil {
		return err

	}
	return nil
}

func (c *Client) GetFunction(fnName string) (*Function, error) {
	log.Debugf("[k8s-GetFunction] namespace=%s fnName=%s", faasNamespace, fnName)

	opts := metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
	}
	deployment, err := c.clientset.AppsV1().Deployments(faasNamespace).Get(context.TODO(), fnName, opts)
	if err != nil {
		return nil, err
	}

	if _, found := deployment.Labels[faasIDLabel]; !found {
		return nil, nil
	}

	return &Function{deployment}, nil
}

func (c *Client) ScaleFunction(fn *Function, replicas int32) error {
	oldReplicas := *fn.Spec.Replicas

	log.Printf("Set replicas - %s %s, %d/%d\n", fn.Name, faasNamespace, replicas, oldReplicas)

	fn.Spec.Replicas = &replicas

	_, err := c.clientset.AppsV1().Deployments(faasNamespace).Update(context.TODO(), fn.Deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
