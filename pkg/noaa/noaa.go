package noaa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bzimmer/transport"
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

// WithHTTPTracing .
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

// WithTransport transport
func WithTransport(t http.RoundTripper) Option {
	return func(c *Client) error {
		if t != nil {
			c.client.Transport = t
		}
		return nil
	}
}

// WithHTTPClient .
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) error {
		if client != nil {
			c.client = client
		}
		return nil
	}
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

// Do executes the request
func (c *Client) Do(req *http.Request, v interface{}) error {
	ctx := req.Context()
	if ctx == nil {
		return errors.New("context must be non-nil")
	}

	res, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		return err
	}
	defer res.Body.Close()

	httpError := res.StatusCode >= http.StatusBadRequest

	var obj interface{}
	if httpError {
		obj = &Fault{}
	} else {
		obj = v
	}

	if obj != nil {
		err := json.NewDecoder(res.Body).Decode(obj)
		if err == io.EOF {
			err = nil // ignore EOF errors caused by empty response body
		}
		if httpError {
			return obj.(error)
		}
		return err
	}

	return nil
}
