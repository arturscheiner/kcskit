package model

type RegistryItem struct {
	ID                 string `json:"id"`
	RegistryName       string `json:"registryName"`
	RegistryType       string `json:"registryType"`
	Description        string `json:"description"`
	RegistryUrl        string `json:"registryUrl"`
	ApiUrl             string `json:"apiUrl"`
	AuthenticationType string `json:"authenticationType"`
	Status             string `json:"status"`
	Message            string `json:"message"`
	LastChecked        string `json:"lastChecked"`
}

type RegistryResponse struct {
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Items []RegistryItem `json:"items"`
}
