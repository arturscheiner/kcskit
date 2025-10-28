package model

import "time"

type CiCdScansListResponse struct {
	Items []CiCdScan `json:"items"`
	Page  int        `json:"page"`
	Total int        `json:"total"`
}

type CiCdScan struct {
	ID           string    `json:"id"`
	ArtifactName string    `json:"artifactName"`
	RiskRating   string    `json:"riskRating"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}
