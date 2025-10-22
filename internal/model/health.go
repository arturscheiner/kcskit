package model

type HealthItem struct {
	ComponentName string `json:"componentName"`
	PodName       string `json:"podName"`
	Status        string `json:"status"`
	Version       string `json:"version"`
	ErrorMessage  string `json:"errorMessage"`
}

type HealthResponse struct {
	Items []HealthItem `json:"items"`
}
