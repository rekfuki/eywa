package types

// Function represents FaaS function
type Function struct {
	Owner     string
	Name      string
	Language  string
	Image     string
	Namespace string
	Labels    map[string]string
}
