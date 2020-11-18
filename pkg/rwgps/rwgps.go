package rwgps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bzimmer/transport"
)

const (
	apiVersion = 2
	baseURL    = "https://ridewithgps.com"
	userAgent  = "(github.com/bzimmer/gravl/rwgps)"
)

// https://ridewithgps.com/api?lang=en

// Client .
type Client struct {
	body   map[string]interface{}
	header http.Header
	client *http.Client

	Users *UsersService
	Trips *TripsService
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
		body:   make(map[string]interface{}),
	}
	// set now, possibly overwritten with options
	c.body["version"] = apiVersion
	// set now, possibly overwritten with options
	c.header.Set("User-Agent", userAgent)
	c.header.Set("Accept", "application/json")
	c.header.Set("Content-type", "application/json")
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// Services used for communicating with RWGPS
	c.Users = &UsersService{client: c}
	c.Trips = &TripsService{client: c}

	return c, nil
}

// WithAuthToken .
func WithAuthToken(authToken string) Option {
	return func(c *Client) error {
		c.body["auth_token"] = authToken
		return nil
	}
}

// WithAPIKey .
func WithAPIKey(apiKey string) Option {
	return func(c *Client) error {
		c.body["apikey"] = apiKey
		return nil
	}
}

// WithAPIVersion .
func WithAPIVersion(version int) Option {
	return func(c *Client) error {
		c.body["version"] = fmt.Sprintf("%d", version)
		return nil
	}
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

// WithAccept .
func WithAccept(accept string) Option {
	return func(c *Client) error {
		if accept != "" {
			c.header.Set("Accept", accept)
		}
		return nil
	}
}

func (c *Client) newBodyReader() (io.Reader, error) {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	err := enc.Encode(c.body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b.Bytes()), nil
}

func (c *Client) newAPIRequest(ctx context.Context, method, uri string) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, uri))
	if err != nil {
		return nil, err
	}
	reader, err := c.newBodyReader()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
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

// Do executes the request
func (c *Client) Do(req *http.Request, v interface{}) error {
	ctx := req.Context()
	res, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
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
