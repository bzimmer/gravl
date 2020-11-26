package visualcrossing

//go:generate go run ../../cmd/genwith/genwith.go --auth --package visualcrossing

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const (
	baseURL   = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/weatherdata"
	userAgent = "(github.com/bzimmer/gravl/pkg/visualcrossing)"
)

// Client .
type Client struct {
	config oauth2.Config
	token  oauth2.Token
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
		token:  oauth2.Token{},
		config: oauth2.Config{},
	}
	// set now, possibly overwritten with options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	c.Forecast = &ForecastService{client: c}

	return c, nil
}

func (c *Client) newAPIRequest(ctx context.Context, method, uri string, values *url.Values) (*http.Request, error) {
	if c.token.AccessToken == "" {
		return nil, errors.New("accessToken required")
	}
	values.Set("key", c.token.AccessToken)
	values.Set("contentType", "json")
	values.Set("locationMode", "array")
	values.Set("shortColumnNames", "false")

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
