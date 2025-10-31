package model

type OllamaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaResponse struct {
	Model           string  `json:"model"`
	CreatedAt       string  `json:"created_at"`
	Message         Message `json:"message"`
	Done            bool    `json:"done"`
	PromptEvalCount int     `json:"prompt_eval_count"`
}

type OllamaHeader struct {
	Command     string
	Cluster     string
	Risk        string
	ReportTitle string
	ApiEndpoint string
}

type TokenizeRequest struct {
	Model   string `json:"model"`
	Content string `json:"content"`
}

type TokenizeResponse struct {
	Tokens []int `json:"tokens"`
}