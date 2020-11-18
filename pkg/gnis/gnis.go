package gnis

import (
	"net/http"

	"github.com/bzimmer/transport"
)

// Client used to communicate with GNIS
type Client struct {
	client *http.Client

	GeoNames *GeoNamesService
}

type service struct {
	client *Client // nolint
}

// Option is used to configure the client
type Option func(*Client) error

// NewClient returns a client ready to query GNIS
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{client: &http.Client{}}
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

// WithTransport option used to configure the RoundTripper
func WithTransport(t http.RoundTripper) Option {
	return func(c *Client) error {
		if t != nil {
			c.client.Transport = t
		}
		return nil
	}
}

// WithHTTPTracing used to configure http tracing
func WithHTTPTracing(debug bool) Option {
	return func(c *Client) error {
		if !debug {
			return nil
		}
		c.client.Transport = &transport.VerboseTransport{
			Transport: c.client.Transport,
		}
		return nil
	}
}
