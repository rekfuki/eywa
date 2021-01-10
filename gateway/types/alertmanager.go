package types

// AlertmanagerAlert represents alert coming from alertmanager
type AlertmanagerAlert struct {
	Status   string  `json:"status"`
	Receiver string  `json:"receiver"`
	Alerts   []alert `json:"alerts"`
}

type alert struct {
	Status      string      `json:"status"`
	Labels      labels      `json:"labels"`
	Annotations annotations `json:"annotations"`
}

type labels struct {
	Severity     string `json:"severity"`
	AlertName    string `json:"alertname"`
	FunctionName string `json:"function_name"`
}

type annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}
