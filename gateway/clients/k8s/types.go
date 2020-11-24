package k8s

import appsv1 "k8s.io/api/apps/v1"

// Function represents FaaS function
type Function struct {
	*appsv1.Deployment
}

// Secret represents k8s secret
type Secret struct {
	Name string
	Data map[string][]byte
}
