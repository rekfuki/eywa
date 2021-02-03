package types

const (
	// UserIDLabel key of the user id label in k8s
	UserIDLabel = "user_id"
	// ImageIDLabel key of the image id label in k8s
	ImageIDLabel = "image_id"
	// SecretIDLabel key of the secret id label in k8s
	SecretIDLabel = "secret_id"
	// FunctionIDLabel key of the function id label in k8s
	FunctionIDLabel = "function_id"
	// UserDefinedNameLabel key of the user defined name label in k8s
	UserDefinedNameLabel = "user_defined_name"

	// LogsSubject is the subject of logs produced to stan
	LogsSubject = "logs"
	// AsyncExecSubject is the subject of asynchronous executions produced to stan
	AsyncExecSubject = "gateway-async"

	// EventHookType represnets event hook type
	EventHookType = 1
	// TimelineHookType represents timeline hook type
	TimelineHookType = 2
)
