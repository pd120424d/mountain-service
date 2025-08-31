package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"go.uber.org/zap"
)

type HTTPClient struct {
	baseURL     string
	httpClient  *http.Client
	serviceAuth auth.ServiceAuth
	logger      utils.Logger
}

type HTTPClientConfig struct {
	BaseURL     string
	Timeout     time.Duration
	ServiceAuth auth.ServiceAuth
	Logger      utils.Logger
}

func NewHTTPClient(config HTTPClientConfig) *HTTPClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &HTTPClient{
		baseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		serviceAuth: config.ServiceAuth,
		logger:      config.Logger,
	}
}

func (c *HTTPClient) Get(ctx context.Context, endpoint string) (*http.Response, error) {
	log := c.logger.WithContext(ctx)
	if strings.Contains(c.baseURL, "/api/v1") && strings.HasPrefix(endpoint, "/api/v1/") {
		log.Warnf("HTTPClient: baseURL already contains /api/v1 and endpoint starts with /api/v1; resulting URL may be double-prefixed: %s + %s", c.baseURL, endpoint)
	}
	return c.doRequest(ctx, "GET", endpoint, nil)
}

func (c *HTTPClient) Post(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	return c.doRequest(ctx, "POST", endpoint, body)
}

func (c *HTTPClient) Put(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	return c.doRequest(ctx, "PUT", endpoint, body)
}

func (c *HTTPClient) Delete(ctx context.Context, endpoint string) (*http.Response, error) {
	return c.doRequest(ctx, "DELETE", endpoint, nil)
}

func (c *HTTPClient) doRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	log := c.logger.WithContext(ctx)
	url := c.baseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Propagate X-Request-ID from context if present
	if rid := utils.RequestIDFromContext(ctx); rid != "" {
		utils.SetRequestIDHeader(req, rid)
	}

	if c.serviceAuth != nil {
		authHeader, err := c.serviceAuth.GetAuthHeader()
		if err != nil {
			return nil, fmt.Errorf("failed to generate auth header: %w", err)
		}
		req.Header.Set("Authorization", authHeader)
	}

	log.Info("Making HTTP request", zap.String("method", method), zap.String("url", url))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error("HTTP request failed", zap.Error(err), zap.String("method", method), zap.String("url", url))
		return nil, fmt.Errorf("request failed: %w", err)
	}

	log.Info("HTTP request completed", zap.String("method", method), zap.String("url", url), zap.Int("status", resp.StatusCode))

	return resp, nil
}
