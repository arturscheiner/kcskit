package controller

import (
	"encoding/json"
	"fmt"

	"github.com/arturscheiner/kcskit/internal/model"
	cfgsvc "github.com/arturscheiner/kcskit/internal/service"
)

// ListClusters calls /v1/clusters with provided rawQuery and returns parsed items, raw body and error.
func ListClusters(cfg model.Config, invalidCert bool, rawQuery string) ([]model.ClusterItem, string, error) {
	client, err := cfgsvc.NewClient(cfg.Endpoint, cfg.Token, invalidCert, cfg.CaCert)
	if err != nil {
		return nil, "", err
	}

	status, body, err := client.Do("GET", "/v1/clusters", rawQuery, nil)
	if err != nil {
		return nil, string(body), err
	}
	if status < 200 || status >= 300 {
		return nil, string(body), fmt.Errorf("received HTTP %d", status)
	}

	var cr model.ClusterResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return nil, string(body), fmt.Errorf("failed to parse clusters JSON: %w", err)
	}
	return cr.Items, string(body), nil
}
