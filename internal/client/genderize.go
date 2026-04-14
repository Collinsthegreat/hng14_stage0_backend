package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Collinsthegreat/hng14_stage0_backend/internal/model"
)

type GenderizeClient interface {
	Predict(ctx context.Context, name string) (*model.GenderizeResponse, error)
}

type genderizeClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewGenderizeClient(httpClient *http.Client, baseURL string) GenderizeClient {
	return &genderizeClient{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *genderizeClient) Predict(ctx context.Context, name string) (*model.GenderizeResponse, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}
	q := u.Query()
	q.Set("name", name)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var data model.GenderizeResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &data, nil
}
