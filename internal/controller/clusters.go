package controller

import (
	"encoding/json"
	"fmt"

	"github.com/arturscheiner/kcskit/internal/model"
	cfgsvc "github.com/arturscheiner/kcskit/internal/service"
)

// ListClusters calls /v1/clusters with provided rawQuery and returns parsed items, raw body and error.
func ListClusters(cfg model.Config, invalidCert bool, rawQuery string) ([]model.ClusterItem, string, string, error) {
	client, err := cfgsvc.NewClient(cfg.Endpoint, cfg.Token, invalidCert, cfg.CaCert)
	if err != nil {
		return nil, "", "", err
	}

	endpoint := "/v1/clusters"
	status, body, err := client.Do("GET", endpoint, rawQuery, nil)
	if err != nil {
		return nil, string(body), endpoint, err
	}
	if status < 200 || status >= 300 {
		return nil, string(body), endpoint, fmt.Errorf("received HTTP %d", status)
	}

	var cr model.ClusterResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return nil, string(body), endpoint, fmt.Errorf("failed to parse clusters JSON: %w", err)
	}
	return cr.Items, string(body), endpoint, nil
}
