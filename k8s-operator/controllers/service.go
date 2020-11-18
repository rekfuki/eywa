package controllers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"

	faasv1alpha1 "eywa/k8s-operator/api/v1alpha1"
	info "eywa/k8s-operator/api/v1alpha1"
)

const functionHTTPPort = 8080

func (r *FunctionReconciler) newService(function *faasv1alpha1.Function) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      function.Spec.Name,
			Namespace: function.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(function, schema.GroupVersionKind{
					Group:   info.GroupVersion.Group,
					Version: info.GroupVersion.Version,
					Kind:    "Function",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: map[string]string{"function": function.Spec.Name},
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: corev1.ProtocolTCP,
					Port:     functionHTTPPort,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(functionHTTPPort),
					},
				},
			},
		},
	}
}
