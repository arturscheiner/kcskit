package service

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/arturscheiner/kcskit/internal/model"
)

// Path returns config file path, ensuring directory exists.
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".kcskit")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "config"), nil
}

// Load reads config YAML into model.Config.
func Load() (model.Config, error) {
	var cfg model.Config
	p, err := Path()
	if err != nil {
		return cfg, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Save merges provided non-empty fields with existing config and writes YAML.
func Save(toSave model.Config) error {
	p, err := Path()
	if err != nil {
		return err
	}

	var existing model.Config
	if b, err := os.ReadFile(p); err == nil {
		_ = yaml.Unmarshal(b, &existing)
	}

	if toSave.Token != "" {
		existing.Token = toSave.Token
	}
	if toSave.Endpoint != "" {
		existing.Endpoint = toSave.Endpoint
	}
	if toSave.CaCert != "" {
		existing.CaCert = toSave.CaCert
	}
	if toSave.AiOllamaEndpoint != "" {
		existing.AiOllamaEndpoint = toSave.AiOllamaEndpoint
	}
	if toSave.AiOllamaModel != "" {
		existing.AiOllamaModel = toSave.AiOllamaModel
	}

	out, err := yaml.Marshal(&existing)
	if err != nil {
		return err
	}
	return os.WriteFile(p, out, 0o600)
}
