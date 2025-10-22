package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/arturscheiner/kcskit/internal/model"
	cfgsvc "github.com/arturscheiner/kcskit/internal/service"
)

// SaveConfig saves token/endpoint/ca_cert (merging with existing).
func SaveConfig(token, endpoint, caCert string) error {
	toSave := model.Config{
		Token:    token,
		Endpoint: endpoint,
		CaCert:   caCert,
	}
	return cfgsvc.Save(toSave)
}

// LoadConfig loads the config file.
func LoadConfig() (model.Config, error) {
	return cfgsvc.Load()
}

// ValidateConfig ensures token and endpoint are present and endpoint looks valid.
func ValidateConfig(cfg model.Config) error {
	if strings.TrimSpace(cfg.Token) == "" {
		return errors.New("token is empty")
	}
	if strings.TrimSpace(cfg.Endpoint) == "" {
		return errors.New("endpoint is empty")
	}
	return nil
}

// TestConfigConnection creates a reusable API client and performs the health action,
// returning the raw response body (JSON) and any error. Caller can parse the JSON.
func TestConfigConnection(cfg model.Config, invalidCert bool) (string, error) {
	client, err := cfgsvc.NewClient(cfg.Endpoint, cfg.Token, invalidCert, cfg.CaCert)
	if err != nil {
		return "", err
	}

	status, body, err := client.Do("GET", "/v1/core-health", "", nil)
	if err != nil {
		return string(body), err
	}
	if status < 200 || status >= 300 {
		return string(body), fmt.Errorf("received HTTP %d", status)
	}
	return string(body), nil
}
