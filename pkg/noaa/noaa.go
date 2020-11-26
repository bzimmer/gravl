package noaa

//go:generate go run ../../cmd/genwith/genwith.go --do --package noaa

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	baseURL   = "https://api.weather.gov"
	userAgent = "(github.com/bzimmer/gravl/pkg/noaa, bzimmer@ziclix.com)"
)

// Client .
type Client struct {
	header http.Header
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
	c := &Client{
		client: &http.Client{},
		header: make(http.Header),
	}
	// set now, possibly overwritten with options
	c.header.Set("User-Agent", userAgent)
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
	for key, values := range c.header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	req.Header.Set("Accept", "application/geo+json")
	return req, nil
}
