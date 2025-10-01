package enricher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GeoClient struct {
	APIKey string
	Client *http.Client
}

type GeoResponse struct {
	CountryName string `json:"country_name"`
	City        string `json:"city"`
}

func NewGeoClient(apiKey string) *GeoClient {
	return &GeoClient{
		APIKey: apiKey,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (g *GeoClient) GetLocation(ctx context.Context, ip string) (*GeoResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(
		"https://api.ipgeolocation.io/ipgeo?apiKey=%s&ip=%s",
		g.APIKey,
		ip,
	), nil)

	if err != nil {
		return nil, err
	}

	resp, err := g.Client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geo API returned %d", resp.StatusCode)
	}

	var data GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
