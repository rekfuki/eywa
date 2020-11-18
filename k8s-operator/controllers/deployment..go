package controllers

import (
	"encoding/json"
	"strings"

	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"

	faasv1alpha1 "eywa/k8s-operator/api/v1alpha1"
	info "eywa/k8s-operator/api/v1alpha1"
)

// newDeployment creates a new Deployment for a Function resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Function resource that 'owns' it.
func (r *FunctionReconciler) newDeployment(
	function *faasv1alpha1.Function,
	existingDeployment *appsv1.Deployment,
	existingSecrets map[string]*corev1.Secret) *appsv1.Deployment {

	envVars := makeEnvVars(function)
	labels := makeLabels(function)
	nodeSelector := makeNodeSelector(function.Spec.Constraints)
	// probes, err := factory.MakeProbes(function)
	// if err != nil {
	// 	glog.Warningf("Function %s probes parsing failed: %v",
	// 		function.Spec.Name, err)
	// }

	resources, err := r.makeResources(function)
	if err != nil {
		log.Warningf("Function %s resources parsing failed: %v",
			function.Spec.Name, err)
	}

	annotations := makeAnnotations(function)

	var serviceAccount string

	if function.Spec.Annotations != nil {
		annotations := *function.Spec.Annotations
		if val, ok := annotations["com.openfaas.serviceaccount"]; ok && len(val) > 0 {
			serviceAccount = val
		}
	}

	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        function.Spec.Name,
			Annotations: annotations,
			Namespace:   function.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(function, schema.GroupVersionKind{
					Group:   info.GroupVersion.Group,
					Version: info.GroupVersion.Version,
					Kind:    "Function",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: r.getReplicas(function, existingDeployment),
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
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":        function.Spec.Name,
					"controller": function.Name,
				},
			},
			RevisionHistoryLimit: int32p(5),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelector,
					Containers: []corev1.Container{
						{
							Name:  function.Spec.Name,
							Image: function.Spec.Image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: int32(8080), Protocol: corev1.ProtocolTCP},
							},
							ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
							Env:             envVars,
							Resources:       *resources,
							// LivenessProbe:   probes.Liveness,
							// ReadinessProbe:  probes.Readiness,
						},
					},
				},
			},
		},
	}

	if len(serviceAccount) > 0 {
		deploymentSpec.Spec.Template.Spec.ServiceAccountName = serviceAccount
	}

	if err := r.UpdateSecrets(function, deploymentSpec, existingSecrets); err != nil {
		log.Warningf("Function %s secrets update failed: %v",
			function.Spec.Name, err)
	}

	return deploymentSpec
}

func makeEnvVars(function *faasv1alpha1.Function) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}

	if len(function.Spec.Handler) > 0 {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "fprocess",
			Value: function.Spec.Handler,
		})
	}

	if function.Spec.Environment != nil {
		for k, v := range *function.Spec.Environment {
			envVars = append(envVars, corev1.EnvVar{
				Name:  k,
				Value: v,
			})
		}
	}

	return envVars
}

func makeLabels(function *faasv1alpha1.Function) map[string]string {
	labels := map[string]string{
		"function":   function.Spec.Name,
		"app":        function.Spec.Name,
		"controller": function.Name,
	}
	if function.Spec.Labels != nil {
		for k, v := range *function.Spec.Labels {
			labels[k] = v
		}
	}

	return labels
}

func makeAnnotations(function *faasv1alpha1.Function) map[string]string {
	annotations := make(map[string]string)

	// copy function annotations
	if function.Spec.Annotations != nil {
		for k, v := range *function.Spec.Annotations {
			annotations[k] = v
		}
	}

	// save function spec in deployment annotations
	// used to detect changes in function spec
	specJSON, err := json.Marshal(function.Spec)
	if err != nil {
		log.Errorf("Failed to marshal function spec: %s", err.Error())
		return annotations
	}

	annotations["com.eywa.rekfuki.dev"] = string(specJSON)
	return annotations
}

func makeNodeSelector(constraints []string) map[string]string {
	selector := make(map[string]string)

	if len(constraints) > 0 {
		for _, constraint := range constraints {
			parts := strings.Split(constraint, "=")

			if len(parts) == 2 {
				selector[parts[0]] = parts[1]
			}
		}
	}

	return selector
}

// deploymentNeedsUpdate determines if the function spec is different from the deployment spec
func deploymentNeedsUpdate(function *faasv1alpha1.Function, deployment *appsv1.Deployment) bool {
	prevFnSpecJson := deployment.ObjectMeta.Annotations["com.eywa.rekfuki.dev"]
	if prevFnSpecJson == "" {
		// is a new deployment or is an old deployment that is missing the annotation
		return true
	}

	prevFnSpec := &faasv1alpha1.FunctionSpec{}
	err := json.Unmarshal([]byte(prevFnSpecJson), prevFnSpec)
	if err != nil {
		log.Errorf("Failed to parse previous function spec: %s", err.Error())
		return true
	}
	prevFn := faasv1alpha1.Function{
		Spec: *prevFnSpec,
	}

	if diff := cmp.Diff(prevFn.Spec, function.Spec); diff != "" {
		log.Infof("Change detected for %s diff\n%s", function.Name, diff)
		return true
	} else {
		log.Infof("No changes detected for %s", function.Name)
	}

	return false
}

func int32p(i int32) *int32 {
	return &i
}

// getReplicas returns the desired number of replicas for a function taking into account
// the min replicas label, HPA, the OF autoscaler and scaled to zero deployments
func (r *FunctionReconciler) getReplicas(function *faasv1alpha1.Function, deployment *appsv1.Deployment) *int32 {
	var minReplicas *int32

	// extract current deployment replicas if specified
	var deploymentReplicas *int32
	if deployment != nil {
		deploymentReplicas = deployment.Spec.Replicas
	}

	// do not set replicas if min replicas is not set
	// and current deployment has no replicas count
	if minReplicas == nil && deploymentReplicas == nil {
		return nil
	}

	// set replicas to min if deployment has no replicas and min replicas exists
	if minReplicas != nil && deploymentReplicas == nil {
		return minReplicas
	}

	// do not override replicas when deployment is scaled to zero
	if deploymentReplicas != nil && *deploymentReplicas == 0 {
		return deploymentReplicas
	}

	// do not override replicas when min is not specified
	if minReplicas == nil && deploymentReplicas != nil {
		return deploymentReplicas
	}

	// do not override HPA or OF autoscaler replicas if the value is greater than min
	if minReplicas != nil && deploymentReplicas != nil {
		if *deploymentReplicas >= *minReplicas {
			return deploymentReplicas
		}
	}

	return minReplicas
}
