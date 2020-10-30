package visualcrossing

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bzimmer/gravl/pkg/common"
)

const (
	visualCrossingURI = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/weatherdata"
	userAgent         = "(github.com/bzimmer/gravl/pkg/visualcrossing)"
)

// Client .
type Client struct {
	apiKey string
	client *http.Client

	Forecast *ForecastService
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

func (c *Client) newAPIRequest(method, uri string, values *url.Values) (*http.Request, error) {
	// these are required
	values.Set("key", c.apiKey)
	values.Set("contentType", "json")
	values.Set("locationMode", "array")

	u, err := url.Parse(fmt.Sprintf("%s/%s?%s", visualCrossingURI, uri, values.Encode()))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Do executes the request
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
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

	var (
		buf    bytes.Buffer
		fault  Fault
		reader = io.TeeReader(res.Body, &buf)
	)

	// VC uses StatusOK for everything, sigh
	if err = json.NewDecoder(reader).Decode(&fault); err != nil {
		return err
	} else if code := fault.ErrorCode; code != 0 && code != http.StatusOK {
		return fault
	}

	if v != nil {
		if err = json.NewDecoder(&buf).Decode(v); err != nil {
			return err
		}
	}

	return nil
}
