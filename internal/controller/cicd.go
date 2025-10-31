package controller

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/arturscheiner/kcskit/internal/model"
	cfgsvc "github.com/arturscheiner/kcskit/internal/service"
)

func ListCicd(cfg model.Config, invalidCert bool, page, limit, sort, by, buildNumber, buildPipeline string) (*model.CiCdScansListResponse, string, string, error) {
	client, err := cfgsvc.NewClient(cfg.Endpoint, cfg.Token, invalidCert, cfg.CaCert)
	if err != nil {
		return nil, "", "", err
	}

	params := url.Values{}
	params.Add("page", page)
	params.Add("limit", limit)
	params.Add("sort", sort)
	params.Add("by", by)
	if buildNumber != "" {
		params.Add("build-number", buildNumber)
	}
	if buildPipeline != "" {
		params.Add("build-pipeline", buildPipeline)
	}

	endpoint := "/v1/scans/ci-cd"
	status, body, err := client.Do("GET", endpoint, params.Encode(), nil)
	if err != nil {
		return nil, string(body), endpoint, err
	}
	if status < 200 || status >= 300 {
		return nil, string(body), endpoint, fmt.Errorf("received HTTP %d", status)
	}

	var cr model.CiCdScansListResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return nil, string(body), endpoint, fmt.Errorf("failed to parse clusters JSON: %w", err)
	}
	return &cr, string(body), endpoint, nil
}