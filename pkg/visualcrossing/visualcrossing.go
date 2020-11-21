package visualcrossing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bzimmer/httpwares"
)

const (
	baseURL   = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/weatherdata"
	userAgent = "(github.com/bzimmer/gravl/pkg/visualcrossing)"
)

// Client .
type Client struct {
	apiKey string
	client *http.Client

	Forecast *ForecastService
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
	}
	// set now, possibly overwritten with options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// Services used for communicating with VisualCrossing
	c.Forecast = &ForecastService{client: c}

	return c, nil
}

// WithAPIKey .
func WithAPIKey(apiKey string) Option {
	return func(c *Client) error {
		c.apiKey = apiKey
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

// WithHTTPTracing .
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

func (c *Client) newAPIRequest(ctx context.Context, method, uri string, values *url.Values) (*http.Request, error) {
	// these are required
	values.Set("key", c.apiKey)
	values.Set("contentType", "json")
	values.Set("locationMode", "array")

	u, err := url.Parse(fmt.Sprintf("%s/%s?%s", baseURL, uri, values.Encode()))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
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

	var (
		buf    bytes.Buffer
		fault  Fault
		reader = io.TeeReader(res.Body, &buf)
	)

	// VC uses StatusOK for everything, sigh
	err = json.NewDecoder(reader).Decode(&fault)
	if err != nil {
		return err
	} else if code := fault.ErrorCode; code != 0 && code != http.StatusOK {
		return &fault
	}

	if v != nil {
		err := json.NewDecoder(&buf).Decode(v)
		if err != nil {
			return err
		}
	}
	return nil
}
