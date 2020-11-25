package cyclinganalytics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"

	"github.com/bzimmer/httpwares"
)

const (
	baseURL   = "https://www.cyclinganalytics.com/api"
	userAgent = "(github.com/bzimmer/gravl/pkg/cyclinganalytics)"
)

// Client .
type Client struct {
	config oauth2.Config
	token  oauth2.Token
	client *http.Client

	Rides *RidesService
}

type service struct {
	client *Client //nolint:golint,structcheck
}

// Option .
type Option func(*Client) error

// Endpoint is CyclingAnalytics's OAuth 2.0 endpoint
var Endpoint = oauth2.Endpoint{
	AuthURL:   fmt.Sprintf("%s/auth", baseURL),
	TokenURL:  fmt.Sprintf("%s/token", baseURL),
	AuthStyle: oauth2.AuthStyleAutoDetect,
}

// NewClient .
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		client: &http.Client{},
		token:  oauth2.Token{},
		config: oauth2.Config{
			Endpoint: Endpoint,
		},
	}
	// set now, possibly overwritten with options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	c.Rides = &RidesService{client: c}

	return c, nil
}

func WithConfig(config oauth2.Config) Option {
	return func(c *Client) error {
		c.config = config
		return nil
	}
}

// WithTokenCredentials provides the tokens for an authenticated user
func WithTokenCredentials(accessToken, refreshToken string, expiry time.Time) Option {
	return func(c *Client) error {
		c.token.AccessToken = accessToken
		c.token.RefreshToken = refreshToken
		c.token.Expiry = expiry
		return nil
	}
}

// WithAPICredentials provides the client api credentials for the application
func WithClientCredentials(clientID, clientSecret string) Option {
	return func(c *Client) error {
		c.config.ClientID = clientID
		c.config.ClientSecret = clientSecret
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
	if c.token.AccessToken == "" {
		return nil, errors.New("accessToken required")
	}
	params := ""
	if values != nil {
		params = values.Encode()
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s?%s", baseURL, uri, params))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", c.token.AccessToken))
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
