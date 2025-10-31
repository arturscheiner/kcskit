package model

type Config struct {
	Token            string `yaml:"token"`
	Endpoint         string `yaml:"endpoint"`
	CaCert           string `yaml:"ca_cert,omitempty"`
	AiOllamaEndpoint string `yaml:"ai_ollama_endpoint,omitempty"`
	AiOllamaModel    string `yaml:"ai_ollama_model,omitempty"`
}
