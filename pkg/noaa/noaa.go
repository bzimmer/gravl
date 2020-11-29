package noaa

//go:generate go run ../../cmd/genwith/genwith.go --do --package noaa

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

type service struct {
	client *Client //nolint:golint,structcheck
}

// Option .
type Option func(*Client) error

// NewClient .
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{client: &http.Client{}}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// Services used for communicating with NOAA
	c.Points = &PointsService{client: c}
	c.GridPoints = &GridPointsService{client: c}

	return c, nil
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
