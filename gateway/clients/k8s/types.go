package k8s

import (
	"time"

	appsv1 "k8s.io/api/apps/v1"
)

// Function represents FaaS function
type Function struct {
	*appsv1.Deployment
}

type FunctionStatus struct {
	Name              string
	Namespace         string
	Image             string
	EnvProcess        string
	Replicas          int32
	MaxReplicas       int32
	MinReplicas       int32
	ScalingFactor     int32
	AvailableReplicas int32
	Annotations       *map[string]string
	Labels            *map[string]string
}

type FunctionZeroScaleResult struct {
	Found     bool
	Available bool
	Duration  time.Duration
}

// Secret represents k8s secret
type Secret struct {
	Name string
	Data map[string][]byte
}
