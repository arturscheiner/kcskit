package model

type ManualJob struct {
	ID           string `json:"id"`
	ScannerName  string `json:"scannerName"`
	Status       string `json:"status"`
	ArtifactName string `json:"artifactName"`
	ArtifactID   string `json:"artifactId"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
