package controller

import (
	"encoding/json"
	"fmt"

	"github.com/arturscheiner/kcskit/internal/model"
	cfgsvc "github.com/arturscheiner/kcskit/internal/service"
)

// ListRegistries fetches image registries via the API and returns parsed items,
// the raw response body and any error. It uses paging/query params suitable for listing.
func ListRegistries(cfg model.Config, invalidCert bool) ([]model.RegistryItem, string, string, error) {
	client, err := cfgsvc.NewClient(cfg.Endpoint, cfg.Token, invalidCert, cfg.CaCert)
	if err != nil {
		return nil, "", "", err
	}

	endpoint := "/v1/registries"
	status, body, err := client.Do("GET", endpoint, "", nil)
	if err != nil {
		return nil, string(body), endpoint, err
	}
	if status != 200 {
		return nil, string(body), endpoint, fmt.Errorf("received HTTP %d", status)
	}

	var items []model.RegistryItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, string(body), endpoint, err
	}
	return items, string(body), endpoint, nil
}
