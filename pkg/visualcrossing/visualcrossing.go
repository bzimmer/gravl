package visualcrossing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// https://www.visualcrossing.com/resources/documentation/weather-api/weather-api-documentation/

const (
	visualCrossingURI = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/weatherdata"
	userAgent         = "(github.com/bzimmer/wta/visualcrossing)"
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
