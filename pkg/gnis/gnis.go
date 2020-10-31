package gnis

import (
	"net/http"

	"github.com/bzimmer/gravl/pkg/common"
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

	// Services used for communicating with GNIS
	c.GeoNames = &GeoNamesService{client: c}

	return c, nil
}

// WithTransport .
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) error {
		if transport != nil {
			c.client.Transport = transport
		}
		return nil
	}
}

// WithVerboseLogging .
func WithVerboseLogging(debug bool) Option {
	return func(c *Client) error {
		if !debug {
			return nil
		}
		c.client.Transport = &common.VerboseTransport{
			Transport: c.client.Transport,
		}
		return nil
	}
}
