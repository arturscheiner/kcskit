package controller

import (
	"encoding/json"
	"fmt"

	"github.com/arturscheiner/kcskit/internal/model"
	cfgsvc "github.com/arturscheiner/kcskit/internal/service"
)

// ListImages calls the /v1/images/registry endpoint with the provided rawQuery (URL-encoded).
// Returns parsed items, raw response body and any error.
func ListImages(cfg model.Config, invalidCert bool, rawQuery string) ([]model.ImageItem, string, error) {
	client, err := cfgsvc.NewClient(cfg.Endpoint, cfg.Token, invalidCert, cfg.CaCert)
	if err != nil {
		return nil, "", err
	}

	status, body, err := client.Do("GET", "/v1/images/registry", rawQuery, nil)
	if err != nil {
		return nil, string(body), err
	}
	if status < 200 || status >= 300 {
		return nil, string(body), fmt.Errorf("received HTTP %d", status)
	}

	var ir model.ImagesResponse
	if err := json.Unmarshal(body, &ir); err != nil {
		return nil, string(body), fmt.Errorf("failed to parse images JSON: %w", err)
	}
	return ir.Items, string(body), nil
}
