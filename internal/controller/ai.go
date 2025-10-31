package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/arturscheiner/kcskit/internal/model"
	"github.com/charmbracelet/glamour"
)

func SendToOllama(jsonOutput string, header model.OllamaHeader) (string, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.AiOllamaEndpoint == "" || cfg.AiOllamaModel == "" {
		return "", fmt.Errorf("ollama endpoint or model not configured")
	}

	endpoint := strings.TrimSuffix(cfg.AiOllamaEndpoint, "/")

	promptContent := fmt.Sprintf("You are an expert on Kaspersky Container Security. You are using a command line utility called kcskit and you have executed the command '%s' that calls the kcs api %s .Evaluate its output and give some insights about this: %s", header.Command, header.ApiEndpoint, jsonOutput)

	var tokenCount int
	tokenizeReqBody := model.TokenizeRequest{
		Model:   cfg.AiOllamaModel,
		Content: promptContent,
	}
	jsonTokenizeBody, err := json.Marshal(tokenizeReqBody)
	if err != nil {
		fmt.Printf("failed to marshal tokenize request body: %v\n", err)
	} else {
		resp, err := http.Post(fmt.Sprintf("%s/api/tokenize", endpoint), "application/json", bytes.NewBuffer(jsonTokenizeBody))
		if err != nil {
			fmt.Printf("failed to send request to ollama for tokenization: %v\n", err)
		} else {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("failed to read tokenize response body: %v\n", err)
			} else {
				var tokenizeResponse model.TokenizeResponse
				if err := json.Unmarshal(body, &tokenizeResponse); err != nil {
					fmt.Printf("failed to unmarshal ollama tokenize response: %v\n", err)
					fmt.Printf("Ollama tokenize response body: %s\n", string(body))
				} else {
					tokenCount = len(tokenizeResponse.Tokens)
				}
			}
		}
	}

	requestBody := model.OllamaRequest{
		Model: cfg.AiOllamaModel,
		Messages: []model.Message{
			{
				Role:    "user",
				Content: promptContent,
			},
		},
		Stream: false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/chat", endpoint), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to send request to ollama: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var ollamaResponse model.OllamaResponse
	if err := json.Unmarshal(body, &ollamaResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal ollama response: %w", err)
	}

	if tokenCount == 0 {
		tokenCount = ollamaResponse.PromptEvalCount
	}

	headerString := fmt.Sprintf(
		"# %s\n\n**Command line:** `%s`\n\n**Date and Time:** %s\n\n**Risk Status Summary:** %s\n\n**Purpose of the Report:** Automated security status and recommendations.\n\n**Model:** %s\n\n**Input Tokens:** %d\n\n---\n\n",
	header.ReportTitle,
	header.Command,
	time.Now().Format(time.RFC1123),
	header.Risk,
	ollamaResponse.Model,
	tokenCount,
	)

	out, err := glamour.Render(headerString+ollamaResponse.Message.Content, "dark")
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	return out, nil
}
