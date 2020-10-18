package noaa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	noaaURI   = "https://api.weather.gov"
	userAgent = "(github.com/bzimmer/wta/noaa, bzimmer@ziclix.com)"
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

	// Services used for talking to NOAA
	c.Points = &PointsService{client: c}
	c.GridPoints = &GridPointsService{client: c}

	return c, nil
}

// RoundTripperFunc wraps a func to make it into a http.RoundTripper. Similar to http.HandleFunc.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip .
func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// WithVerboseLogging .
func WithVerboseLogging(debug bool) func(*Client) error {
	return func(client *Client) error {
		if !debug {
			return nil
		}
		transport := client.client.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}
		client.client.Transport = RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			dump, _ := httputil.DumpRequestOut(req, true)
			log.Debug().Str("req", string(dump)).Msg("sending")
			res, err := transport.RoundTrip(req)
			dump, _ = httputil.DumpResponse(res, true)
			log.Debug().Str("res", string(dump)).Msg("received")
			return res, err
		})
		return nil
	}
}

// WithTimeout timeout
func WithTimeout(timeout time.Duration) func(*Client) error {
	return func(client *Client) error {
		client.client.Timeout = timeout
		return nil
	}
}

// WithHTTPClient .
func WithHTTPClient(client *http.Client) func(c *Client) error {
	return func(c *Client) error {
		if client != nil {
			c.client = client
		}
		return nil
	}
}

// WithAccept .
func WithAccept(accept string) func(c *Client) error {
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

	if v != nil {
		err := json.NewDecoder(res.Body).Decode(v)
		if err == io.EOF {
			err = nil // ignore EOF errors caused by empty response body
		}
		return err
	}

	return nil
}
