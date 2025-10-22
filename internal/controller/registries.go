package controller

import (
	"encoding/json"
	"fmt"

	"github.com/arturscheiner/kcskit/internal/model"
	cfgsvc "github.com/arturscheiner/kcskit/internal/service"
)

// ListRegistries fetches image registries via the API and returns parsed items,
// the raw response body and any error. It uses paging/query params suitable for listing.
func ListRegistries(cfg model.Config, invalidCert bool) ([]model.RegistryItem, string, error) {
	// pass cfg.CaCert to NewClient (signature: baseURL, token, invalidCert, caCert)
	client, err := cfgsvc.NewClient(cfg.Endpoint, cfg.Token, invalidCert, cfg.CaCert)
	if err != nil {
		return nil, "", err
	}

	rawQuery := "page=1&limit=50&sort=name&by=asc"
	status, body, err := client.Do("GET", "/v1/integrations/image-registries", rawQuery, nil)
	if err != nil {
		return nil, string(body), err
	}
	if status < 200 || status >= 300 {
		return nil, string(body), fmt.Errorf("received HTTP %d", status)
	}

	var rr model.RegistryResponse
	if err := json.Unmarshal(body, &rr); err != nil {
		return nil, string(body), fmt.Errorf("failed to parse registries JSON: %w", err)
	}
	return rr.Items, string(body), nil
}
