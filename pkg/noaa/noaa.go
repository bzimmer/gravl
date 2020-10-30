package noaa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/bzimmer/gravl/pkg/common"
)

const (
	noaaURI   = "https://api.weather.gov"
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
	client *Client
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
	c.header.Set("Accept", "application/geo+json")
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

// WithTransport transport
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) error {
		if transport != nil {
			c.client.Transport = transport
		}
		return nil
	}
}

// WithTimeout timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		c.client.Timeout = timeout
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

// WithAccept .
func WithAccept(accept string) Option {
	return func(c *Client) error {
		if accept != "" {
			c.header.Set("Accept", accept)
		}
		return nil
	}
}

func (c *Client) newAPIRequest(method, uri string) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", noaaURI, uri))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	for key, values := range c.header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return req, nil
}

func withContext(ctx context.Context, req *http.Request) *http.Request {
	// No-op
	return req
}

// Do executes the request
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
	if ctx == nil {
		return errors.New("context must be non-nil")
	}
	req = withContext(ctx, req)

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
