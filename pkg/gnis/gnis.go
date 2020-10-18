package gnis

import (
	"net/http"
)

const (
	gnisLength = 20
	userAgent  = "github.com/bzimmer/wta"
	baseURL    = "https://geonames.usgs.gov/docs/stategaz/%s_Features.zip"
)

// Client .
type Client struct {
	client *http.Client

	GeoNames *GeoNamesService
}

type service struct {
	client *Client
}

// Option .
type Option func(*Client) error

// NewClient .
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		client: &http.Client{},
	}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// Services used for talking to the GNIS website
	c.GeoNames = &GeoNamesService{client: c}

	return c, nil
}
