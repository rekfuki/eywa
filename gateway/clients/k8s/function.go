package k8s

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// DeployFunction deploys function to the k8s cluster
func (c *Client) DeployFunction(request *DeployFunctionRequest) (*FunctionStatus, error) {
	request.Limits = &FunctionResources{
		CPU:    c.limitRange.MaxCPU,
		Memory: c.limitRange.MaxMem,
	}

	deployment, err := c.buildDeployment(request)
	if err != nil {
		return nil, err
	}

	setSecrets(request, deployment)

	deployment, err = c.clientset.AppsV1().
		Deployments(faasNamespace).
		Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	service := buildService(request)
	_, err = c.clientset.CoreV1().
		Services(faasNamespace).
		Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return deploymentToFunction(deployment)
}

// UpdateFunction updates function deployment
func (c *Client) UpdateFunction(oldName string, request *DeployFunctionRequest) (*FunctionStatus, error) {
	context := context.TODO()
	deployment, err := c.clientset.AppsV1().
		Deployments(faasNamespace).
		Get(context, oldName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		request.Limits = &FunctionResources{
			CPU:    c.limitRange.MaxCPU,
			Memory: c.limitRange.MaxMem,
		}

		baseDeployment, err := c.buildDeployment(request)
		if err != nil {
			return nil, err
		}

		deployment.Spec.Template.Spec.Containers[0].Env = baseDeployment.Spec.Template.Spec.Containers[0].Env
		deployment.Spec.Template.Spec.Containers[0].Image = baseDeployment.Spec.Template.Spec.Containers[0].Image

		deployment.Spec.Template.Spec.NodeSelector = baseDeployment.Spec.Template.Spec.NodeSelector

		deployment.Spec.Template.ObjectMeta.Labels = baseDeployment.Spec.Template.ObjectMeta.Labels
		deployment.ObjectMeta.Labels = baseDeployment.ObjectMeta.Labels

		deployment.Annotations = baseDeployment.Annotations
		deployment.Spec.Template.Annotations = baseDeployment.Spec.Template.Annotations
		deployment.Spec.Template.ObjectMeta.Annotations = baseDeployment.Spec.Template.ObjectMeta.Annotations

		deployment.Spec.Template.Spec.Containers[0].Resources = baseDeployment.Spec.Template.Spec.Containers[0].Resources

		setSecrets(request, deployment)

		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = baseDeployment.Spec.Template.Spec.Containers[0].LivenessProbe
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = baseDeployment.Spec.Template.Spec.Containers[0].ReadinessProbe

		deployment.Spec.Replicas = baseDeployment.Spec.Replicas
	}

	// This might cause side effects.

	deployment, err = c.clientset.AppsV1().
		Deployments(faasNamespace).
		Update(context, deployment, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return deploymentToFunction(deployment)
}

func buildService(request *DeployFunctionRequest) *corev1.Service {
	annotations := map[string]string{}
	if len(request.Annotations) > 0 {
		annotations = request.Annotations
	}

	// DNS records cannot start with a number.
	// Since all the names are UUIDs, there is a
	// high chance that a number will be the first character.
	serviceName := "s-" + request.Service
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        serviceName,
			Annotations: annotations,
			Labels:      request.Labels,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				faasIDLabel: request.Service, // Deployment can start with a number
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
}

func (c *Client) buildDeployment(request *DeployFunctionRequest) (*appsv1.Deployment, error) {
	envVars := []corev1.EnvVar{{
		Name:  "mongodb_host",
		Value: c.mongoDBHost,
	}}
	for k, v := range request.EnvVars {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	request.Labels[updatedAtLabel] = fmt.Sprint(time.Now().Unix())
	labels := map[string]string{
		faasIDLabel: request.Service,
	}

	for k, v := range request.Labels {
		labels[k] = v
	}

	if request.MinReplicas == 0 {
		request.MinReplicas = defaultMinReplicas
	}

	if request.MaxReplicas == 0 {
		request.MaxReplicas = defaultMaxReplicas
	}

	if request.ScalingFactor == 0 {
		request.ScalingFactor = defaultScalingFactor
	}

	labels[faasIDLabel] = request.Service
	labels[faasMinReplicasIDLabel] = strconv.Itoa(request.MinReplicas)
	labels[faasMaxReplicasIDLabel] = strconv.Itoa(request.MaxReplicas)
	labels[faasScaleFactorIDLabel] = strconv.Itoa(request.ScalingFactor)

	replicaCount := int32(request.MinReplicas)
	resources := &apiv1.ResourceRequirements{
		Limits:   apiv1.ResourceList{},
		Requests: apiv1.ResourceList{},
	}

	// Set Memory limits
	if request.Limits != nil && len(request.Limits.Memory) > 0 {
		qty, err := resource.ParseQuantity(request.Limits.Memory)
		if err != nil {
			return nil, err
		}
		resources.Limits[apiv1.ResourceMemory] = qty
	}

	if request.Requests != nil && len(request.Requests.Memory) > 0 {
		qty, err := resource.ParseQuantity(request.Requests.Memory)
		if err != nil {
			return nil, err
		}
		resources.Requests[apiv1.ResourceMemory] = qty
	}

	// Set CPU limits
	if request.Limits != nil && len(request.Limits.CPU) > 0 {
		qty, err := resource.ParseQuantity(request.Limits.CPU)
		if err != nil {
			return nil, err
		}
		resources.Limits[apiv1.ResourceCPU] = qty
	}

	if request.Requests != nil && len(request.Requests.CPU) > 0 {
		qty, err := resource.ParseQuantity(request.Requests.CPU)
		if err != nil {
			return nil, err
		}
		resources.Requests[apiv1.ResourceCPU] = qty
	}

	imagePullPolicy := apiv1.PullAlways

	annotations := map[string]string{}
	if len(request.Annotations) > 0 {
		annotations = request.Annotations
	}

	var handler corev1.Handler
	initialDelaySeconds := initialDelaySeconds

	handler = corev1.Handler{
		HTTPGet: &corev1.HTTPGetAction{
			Path: probePathValue,
			Port: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(8080),
			},
		},
	}

	probe := &corev1.Probe{
		Handler:             handler,
		InitialDelaySeconds: int32(initialDelaySeconds),
		TimeoutSeconds:      int32(timeoutSeconds),
		PeriodSeconds:       int32(periodSeconds),
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        request.Service,
			Annotations: annotations,
			Labels:      labels,
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
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: dockerPullSecret,
						},
					},
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
							LivenessProbe:   probe,
							ReadinessProbe:  probe,
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
					DNSPolicy:     corev1.DNSClusterFirst,
				},
			},
		},
	}

	return deployment, nil
}

func setSecrets(request *DeployFunctionRequest, deployment *appsv1.Deployment) {
	secretVolumeProjections := []apiv1.VolumeProjection{}
	for _, secret := range request.Secrets {
		projectedPaths := []apiv1.KeyToPath{}
		for secretKey := range secret.Data {
			projectedPaths = append(projectedPaths, apiv1.KeyToPath{Key: secretKey, Path: secret.Name + "/" + secretKey})
		}

		projection := &apiv1.SecretProjection{Items: projectedPaths}
		projection.Name = secret.Name
		secretProjection := apiv1.VolumeProjection{Secret: projection}
		secretVolumeProjections = append(secretVolumeProjections, secretProjection)
	}

	volumeName := fmt.Sprintf("%s-projected-secrets", request.Service)
	projectedSecrets := apiv1.Volume{
		Name: volumeName,
		VolumeSource: apiv1.VolumeSource{
			Projected: &apiv1.ProjectedVolumeSource{
				Sources: secretVolumeProjections,
			},
		},
	}

	existingVolumes := removeVolume(volumeName, deployment.Spec.Template.Spec.Volumes)
	deployment.Spec.Template.Spec.Volumes = existingVolumes
	if len(secretVolumeProjections) > 0 {
		deployment.Spec.Template.Spec.Volumes = append(existingVolumes, projectedSecrets)
	}

	updatedContainers := []apiv1.Container{}
	for _, container := range deployment.Spec.Template.Spec.Containers {
		mount := apiv1.VolumeMount{
			Name:      volumeName,
			ReadOnly:  true,
			MountPath: faasSecretMount,
		}

		container.VolumeMounts = removeVolumeMount(volumeName, container.VolumeMounts)
		if len(secretVolumeProjections) > 0 {
			container.VolumeMounts = append(container.VolumeMounts, mount)
		}

		updatedContainers = append(updatedContainers, container)
	}

	deployment.Spec.Template.Spec.Containers = updatedContainers
}

func removeVolume(volumeName string, volumes []corev1.Volume) []corev1.Volume {
	newVolumes := volumes[:0]
	for _, v := range volumes {
		if v.Name != volumeName {
			newVolumes = append(newVolumes, v)
		}
	}

	return newVolumes
}

func removeVolumeMount(volumeName string, mounts []corev1.VolumeMount) []corev1.VolumeMount {
	newMounts := mounts[:0]
	for _, v := range mounts {
		if v.Name != volumeName {
			newMounts = append(newMounts, v)
		}
	}

	return newMounts
}

// DeleteFunction deletes the function deployment and service
func (c *Client) DeleteFunction(fnName string) error {
	foregroundPolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &foregroundPolicy}

	if err := c.clientset.AppsV1().
		Deployments(faasNamespace).
		Delete(context.TODO(), fnName, *opts); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	serviceName := "s-" + fnName
	if err := c.clientset.CoreV1().
		Services(faasNamespace).
		Delete(context.TODO(), serviceName, *opts); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

// GetFunctionStatus returns status of the function from k8s
func (c *Client) GetFunctionStatus(filter Selector) (*FunctionStatus, error) {
	return c.getFunctionStatus(filter)
}

// GetFunctionStatusFiltered returns function filtered by selector
func (c *Client) GetFunctionStatusFiltered(filter Selector) (*FunctionStatus, error) {
	return c.getFunctionStatus(filter)
}

func (c *Client) getFunctionStatus(filter Selector) (*FunctionStatus, error) {
	deployments, err := c.listDeployments(filter)
	if err != nil {
		return nil, err
	}

	if len(deployments.Items) == 0 {
		return nil, nil
	}

	if len(deployments.Items) != 1 {
		log.Warnf("K8s returned more than one result when only one was expected: %#v", filter)
	}

	return deploymentToFunction(&deployments.Items[0])
}

// GetFunctionsStatus returns all functions with faas id label
func (c *Client) GetFunctionsStatus() ([]FunctionStatus, error) {
	return c.getFunctionsStatus(LabelSelector().Exists(faasIDLabel))
}

// GetFunctionsStatusFiltered returns functions filtered by userID
func (c *Client) GetFunctionsStatusFiltered(filter Selector) ([]FunctionStatus, error) {
	return c.getFunctionsStatus(filter)
}

func (c *Client) getFunctionsStatus(filter Selector) ([]FunctionStatus, error) {
	functions, err := c.listDeployments(filter)
	if err != nil {
		return nil, err
	}

	fs := []FunctionStatus{}
	for _, d := range functions.Items {
		f, err := deploymentToFunction(&d)
		if err != nil {
			return nil, err
		}
		fs = append(fs, *f)
	}

	return fs, nil
}

func (c *Client) listDeployments(filter Selector) (*appsv1.DeploymentList, error) {
	requirement, err := labels.NewRequirement(faasIDLabel, selection.Exists, []string{})
	if err != nil {
		return nil, err
	}
	filter.Add(*requirement)

	opts := metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		LabelSelector: filter.String(),
	}

	return c.clientset.AppsV1().Deployments(faasNamespace).List(context.TODO(), opts)
}

func deploymentToFunction(deployment *appsv1.Deployment) (*FunctionStatus, error) {
	var replicas int
	if deployment.Spec.Replicas != nil {
		replicas = int(*deployment.Spec.Replicas)
	}

	function := &FunctionStatus{
		Name:              deployment.Name,
		Replicas:          replicas,
		Image:             deployment.Spec.Template.Spec.Containers[0].Image,
		AvailableReplicas: int(deployment.Status.AvailableReplicas),
		Labels:            deployment.Spec.Template.Labels,
		Annotations:       deployment.Spec.Template.Annotations,
		Namespace:         deployment.Namespace,
		CreatedAt:         deployment.ObjectMeta.CreationTimestamp.Time,
		Env:               map[string]string{},
	}

	if deployment.ObjectMeta.DeletionTimestamp != nil {
		function.DeletedAt = &deployment.ObjectMeta.DeletionTimestamp.Time
	}

	for _, v := range deployment.Spec.Template.Spec.Volumes {
		if v.Projected != nil {
			for _, s := range v.Projected.Sources {
				if s.Secret != nil {
					function.MountedSecrets = append(function.MountedSecrets, s.Secret.Name)
				}
			}
		}
	}

	c := deployment.Spec.Template.Spec.Containers[0]
	for _, v := range c.Env {
		function.Env[v.Name] = v.Value
	}

	function.Limits = &FunctionResources{
		CPU:    c.Resources.Limits.Cpu().String(),
		Memory: c.Resources.Limits.Memory().String(),
	}

	function.Requests = &FunctionResources{
		CPU:    c.Resources.Requests.Cpu().String(),
		Memory: c.Resources.Requests.Memory().String(),
	}

	for k, v := range deployment.Spec.Template.Labels {
		switch k {
		case faasMinReplicasIDLabel, faasMaxReplicasIDLabel, faasScaleFactorIDLabel:
			i64, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return nil, err
			}

			i := int(i64)

			if k == faasMinReplicasIDLabel {
				function.MinReplicas = i
			} else if k == faasMaxReplicasIDLabel {
				function.MaxReplicas = i
			} else {
				function.ScalingFactor = i
			}
		case updatedAtLabel:
			t, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			function.UpdatedAt = time.Unix(t, 0)
		}
	}

	function.Available = true
	if deployment.Status.ReadyReplicas == 0 {
		if function.MinReplicas > 0 || deployment.Status.UnavailableReplicas > 0 {
			function.Available = false
		}
	}

	return function, nil
}

// ScaleFunction scales the function to specified replicas
func (c *Client) ScaleFunction(filter Selector, replicas int) error {
	deployments, err := c.listDeployments(filter)
	if err != nil {
		return err
	}

	if len(deployments.Items) == 0 {
		return fmt.Errorf("Failed to scale. Function %q not found: %s", filter.String(), err)
	}

	if len(deployments.Items) > 1 {
		log.Warnf("Filter %q matched more than one function, when only one was expected", filter.String())
	}

	deployment := &deployments.Items[0]
	oldReplicas := *deployment.Spec.Replicas

	log.Printf("Set replicas - %s %s, %d/%d\n", deployment.Name, faasNamespace, replicas, oldReplicas)

	i32Replicas := int32(replicas)
	deployment.Spec.Replicas = &i32Replicas

	_, err = c.clientset.AppsV1().
		Deployments(faasNamespace).
		Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// ScaleFromZero scales the function from zero replicas to desired
func (c *Client) ScaleFromZero(filter Selector) (*FunctionZeroScaleResult, error) {
	start := time.Now()

	if val, found := c.cache.Get(filter.String()); found {
		cached, ok := val.(*FunctionStatus)
		if !ok {
			return &FunctionZeroScaleResult{
				Available: false,
				Found:     false,
			}, fmt.Errorf("Cache error: expected %T, received %T", &FunctionStatus{}, val)
		}
		if cached.AvailableReplicas > 0 {
			return &FunctionZeroScaleResult{
				Available:      true,
				Found:          true,
				FunctionStatus: cached,
			}, nil
		}
	}

	functionStatus, err := c.GetFunctionStatus(filter)
	if functionStatus == nil || err != nil {
		return &FunctionZeroScaleResult{
			Available: false,
			Found:     false,
			Duration:  time.Since(start),
		}, err
	}

	c.cache.Set(filter.String(), functionStatus, cache.DefaultExpiration)

	if functionStatus.AvailableReplicas == 0 {
		minReplicas := 1
		if functionStatus.MinReplicas > 0 {
			minReplicas = functionStatus.MinReplicas
		}

		// TODO: move to config
		attempts := 20
		interval := time.Millisecond * 50

		err := backoff(func(attempt int) error {
			functionStatus, err := c.GetFunctionStatus(filter)
			if err != nil {
				return err
			}

			c.cache.Set(filter.String(), functionStatus, cache.DefaultExpiration)

			if functionStatus.AvailableReplicas > 0 {
				return nil
			}

			err = c.ScaleFunction(filter, minReplicas)
			if err != nil {
				return fmt.Errorf("Failed to scale function %q, err: %s", filter.String(), err)
			}
			return nil
		}, attempts, interval)

		if err != nil {
			return &FunctionZeroScaleResult{
				Available: false,
				Found:     true,
				Duration:  time.Since(start),
			}, err
		}

		// TODO: move to config
		maxPollCount := 1000

		for i := 0; i < maxPollCount; i++ {
			functionStatus, err := c.GetFunctionStatus(filter)
			if err != nil || functionStatus == nil {
				return &FunctionZeroScaleResult{
					Available: false,
					Found:     true,
					Duration:  time.Since(start),
				}, err
			}

			c.cache.Set(filter.String(), functionStatus, cache.DefaultExpiration)
			totalTime := time.Since(start)

			if functionStatus.AvailableReplicas > 0 {
				log.Printf("Function %q scaled successfully in %fs. Available replicas: %d",
					filter.String(), totalTime.Seconds(), functionStatus.AvailableReplicas)

				return &FunctionZeroScaleResult{
					Available:      true,
					Found:          true,
					Duration:       totalTime,
					FunctionStatus: functionStatus,
				}, nil
			}

			time.Sleep(interval)
		}
	}

	return &FunctionZeroScaleResult{
		Available:      true,
		Found:          true,
		Duration:       time.Since(start),
		FunctionStatus: functionStatus,
	}, nil
}

type routine func(attempt int) error

func backoff(r routine, attempts int, interval time.Duration) error {
	var err error

	for i := 0; i < attempts; i++ {
		res := r(i)
		if res != nil {
			err = res

			log.Printf("Attempt: %d, had error: %s\n", i, res)
		} else {
			err = nil
			break
		}
		time.Sleep(interval)
	}
	return err
}
