package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FunctionSpec defines the desired state of Function
type FunctionSpec struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	// +optional
	Handler string `json:"handler,omitempty"`
	// +optional
	Annotations *map[string]string `json:"annotations,omitempty"`
	// +optional
	Labels *map[string]string `json:"labels,omitempty"`
	// +optional
	Environment *map[string]string `json:"environment,omitempty"`
	// +optional
	Constraints []string `json:"constraints,omitempty"`
	// +optional
	Secrets []string `json:"secrets,omitempty"`
	// +optional
	Limits *FunctionResources `json:"limits,omitempty"`
	// +optional
	Requests *FunctionResources `json:"requests,omitempty"`
}

// FunctionResources is used to set CPU and memory limits and requests
type FunctionResources struct {
	Memory string `json:"memory,omitempty"`
	CPU    string `json:"cpu,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Function is the Schema for the functions API
type Function struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FunctionSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// FunctionList contains a list of Function
type FunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Function `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Function{}, &FunctionList{})
}
