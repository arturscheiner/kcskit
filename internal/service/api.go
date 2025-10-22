package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// APIClient is a reusable client bound to a base URL and token.
type APIClient struct {
	BaseURL *url.URL
	Token   string
	HTTP    *http.Client
}

// NewClient creates an APIClient from baseURL and token.
// If invalidCert is true TLS verification is skipped.
// If caCert PEM text is provided and invalidCert is false, it will be used as RootCAs.
func NewClient(baseURL, token string, invalidCert bool, caCert string) (*APIClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid base URL: %s", baseURL)
	}

	tr := &http.Transport{}
	tlsCfg := &tls.Config{}

	if invalidCert {
		tlsCfg.InsecureSkipVerify = true //nolint:gosec
	} else if strings.TrimSpace(caCert) != "" {
		pool := x509.NewCertPool()
		ok := pool.AppendCertsFromPEM([]byte(caCert))
		if !ok {
			return nil, fmt.Errorf("failed to parse provided ca_cert PEM")
		}
		tlsCfg.RootCAs = pool
		// Note: we leave InsecureSkipVerify false so verification uses provided RootCAs.
	}
	// only set TLSClientConfig if we modified tlsCfg (to avoid nil case)
	if tlsCfg.InsecureSkipVerify || tlsCfg.RootCAs != nil {
		tr.TLSClientConfig = tlsCfg
	}

	client := &http.Client{
		Timeout:   8 * time.Second,
		Transport: tr,
	}
	return &APIClient{
		BaseURL: u,
		Token:   token,
		HTTP:    client,
	}, nil
}

// Do performs an HTTP request to actionPath (relative to base URL).
// method: "GET", "POST", ...
// actionPath: e.g. "/v1/core-health" or "v1/core-health"
// rawQuery: optional query string (without leading '?') — can be empty
// headers: optional additional headers
// Returns status code, response body bytes and error.
func (c *APIClient) Do(method, actionPath, rawQuery string, headers map[string]string) (int, []byte, error) {
	actionPath = strings.TrimSpace(actionPath)
	if i := strings.Index(actionPath, "?"); i != -1 && rawQuery == "" {
		rawQuery = actionPath[i+1:]
		actionPath = actionPath[:i]
	}
	joined := path.Join(strings.TrimSuffix(c.BaseURL.Path, "/"), strings.TrimPrefix(actionPath, "/"))
	if !strings.HasPrefix(joined, "/") {
		joined = "/" + joined
	}
	u := c.BaseURL.ResolveReference(&url.URL{Path: joined, RawQuery: rawQuery})

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("accept", "application/json")
	if c.Token != "" {
		req.Header.Set("Tron-Token", c.Token)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, b, nil
}

// PostJSON marshals payload to JSON and sends it as a POST request to actionPath.
// It returns HTTP status, response body bytes and error.
func (c *APIClient) PostJSON(actionPath, rawQuery string, payload interface{}, headers map[string]string) (int, []byte, error) {
	// marshal payload
	b, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, err
	}

	// build URL similar to Do()
	actionPath = strings.TrimSpace(actionPath)
	if i := strings.Index(actionPath, "?"); i != -1 && rawQuery == "" {
		rawQuery = actionPath[i+1:]
		actionPath = actionPath[:i]
	}
	joined := path.Join(strings.TrimSuffix(c.BaseURL.Path, "/"), strings.TrimPrefix(actionPath, "/"))
	if !strings.HasPrefix(joined, "/") {
		joined = "/" + joined
	}
	u := c.BaseURL.ResolveReference(&url.URL{Path: joined, RawQuery: rawQuery})

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(b))
	if err != nil {
		return 0, nil, err
	}
	// default headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Tron-Token", c.Token)
	}
	// extra headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, respBody, nil
}
