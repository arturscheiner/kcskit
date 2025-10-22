package controller

import (
	"encoding/json"
	"fmt"

	"github.com/arturscheiner/kcskit/internal/model"
	cfgsvc "github.com/arturscheiner/kcskit/internal/service"
)

// CreateScan triggers a manual scan for an artifact in a registry.
// Returns parsed ManualJob, raw response body and error.
func CreateScan(cfg model.Config, invalidCert bool, artifact string, registryID string) (model.ManualJob, string, error) {
	var job model.ManualJob

	client, err := cfgsvc.NewClient(cfg.Endpoint, cfg.Token, invalidCert, cfg.CaCert)
	if err != nil {
		return job, "", err
	}

	payload := map[string]string{
		"artifact":   artifact,
		"registryId": registryID,
	}

	status, body, err := client.PostJSON("/v1/scans", "", payload, nil)
	if err != nil {
		return job, string(body), err
	}
	if status < 200 || status >= 300 {
		return job, string(body), fmt.Errorf("received HTTP %d", status)
	}

	if err := json.Unmarshal(body, &job); err != nil {
		return job, string(body), fmt.Errorf("failed to parse scan response: %w", err)
	}
	return job, string(body), nil
}
