package model

type ImageItem struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	ImageRegistryName string `json:"imageRegistryName"`
	NonCompliant      int    `json:"nonCompliant"`
	Total             int    `json:"total"`
	Errors            int    `json:"errors"`
	Process           int    `json:"process"`
	RiskRating        string `json:"riskRating"`
	Public            bool   `json:"public"`
}

type ImagesResponse struct {
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Items []ImageItem `json:"items"`
}
