package model

type ClusterItem struct {
	ID           string `json:"id"`
	AgentGroupId string `json:"agentGroupId"`
	ClusterName  string `json:"clusterName"`
	Orchestrator string `json:"orchestrator"`
	Namespaces   int    `json:"namespaces"`
	RiskRating   string `json:"riskRating"`
}

type ClusterResponse struct {
	Total int           `json:"total"`
	Page  int           `json:"page"`
	Items []ClusterItem `json:"items"`
}
