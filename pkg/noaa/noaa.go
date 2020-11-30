package noaa

//go:generate go run ../../cmd/genwith/genwith.go --do --client --package noaa

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bzimmer/gravl/pkg"
)

const (
	baseURL = "https://api.weather.gov"
)

// Client .
type Client struct {
	client *http.Client

	Points     *PointsService
	GridPoints *GridPointsService
}

func withServices(c *Client) {
	c.Points = &PointsService{client: c}
	c.GridPoints = &GridPointsService{client: c}
}

func (c *Client) newAPIRequest(ctx context.Context, method, uri string) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, uri))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", pkg.UserAgent)
	req.Header.Set("Accept", "application/geo+json")
	return req, nil
}
