// Code generated by "genwith.go --client --ratelimit --package web"; DO NOT EDIT.

package web

import (
	"errors"
	"net/http"

	"github.com/bzimmer/httpwares"
	"golang.org/x/time/rate"
)

type service struct {
	client *Client //nolint:golint,structcheck
}

// Option provides a configuration mechanism for a Client
type Option func(*Client) error

// NewClient creates a new client and applies all provided Options
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		client: &http.Client{},
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if err := withServices(c); err != nil {
		return nil, err
	}
	return c, nil
}

// WithRateLimiter rate limits the client's api calls
func WithRateLimiter(r *rate.Limiter) Option {
	return func(c *Client) error {
		if r == nil {
			return errors.New("nil limiter")
		}
		c.client.Transport = &httpwares.RateLimitTransport{
			Limiter:   r,
			Transport: c.client.Transport,
		}
		return nil
	}
}

// WithHTTPTracing enables tracing http calls.
func WithHTTPTracing(debug bool) Option {
	return func(c *Client) error {
		if !debug {
			return nil
		}
		c.client.Transport = &httpwares.VerboseTransport{
			Transport: c.client.Transport,
		}
		return nil
	}
}

// WithTransport sets the underlying http client transport.
func WithTransport(t http.RoundTripper) Option {
	return func(c *Client) error {
		if t == nil {
			return errors.New("nil transport")
		}
		c.client.Transport = t
		return nil
	}
}

// WithHTTPClient sets the underlying http client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) error {
		if client == nil {
			return errors.New("nil client")
		}
		c.client = client
		return nil
	}
}
