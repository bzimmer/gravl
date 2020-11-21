package strava

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/bzimmer/httpwares"
	"github.com/markbates/goth"
)

const (
	baseURL = "https://www.strava.com/api/v3"
)

// Client client
type Client struct {
	stravaKey    string
	stravaSecret string
	accessToken  string
	refreshToken string

	provider goth.Provider
	client   *http.Client

	Auth     *AuthService
	Route    *RouteService
	Webhook  *WebhookService
	Athlete  *AthleteService
	Activity *ActivityService
}

type service struct {
	client *Client //nolint:golint,structcheck
}

// Option .
type Option func(*Client) error

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

// WithTransport configures the http client to use the provided transport
func WithTransport(t http.RoundTripper) Option {
	return func(c *Client) error {
		if t != nil {
			c.client.Transport = t
		}
		return nil
	}
}

// WithWebhookCredentials provides the Strava credentials
func WithWebhookCredentials(stravaKey, stravaSecret string) Option {
	return func(c *Client) error {
		c.stravaKey = stravaKey
		c.stravaSecret = stravaSecret
		return nil
	}
}

// WithAPICredentials provides the Strava credentials
func WithAPICredentials(accessToken, refreshToken string) Option {
	return func(c *Client) error {
		c.accessToken = accessToken
		c.refreshToken = refreshToken
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

// WithProvider .
func WithProvider(provider goth.Provider) Option {
	return func(c *Client) error {
		c.provider = provider
		return nil
	}
}

// NewClient creates new clients
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{client: &http.Client{}}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// Services used for talking to Strava
	c.Auth = &AuthService{client: c}
	c.Route = &RouteService{client: c}
	c.Webhook = &WebhookService{client: c}
	c.Athlete = &AthleteService{client: c}
	c.Activity = &ActivityService{client: c}

	return c, nil
}

func (c *Client) newAPIRequest(ctx context.Context, method, uri string) (*http.Request, error) { // nolint:unparam
	if c.accessToken == "" {
		return nil, errors.New("accessToken required")
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, uri))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", c.accessToken))
	return req, nil
}

func (c *Client) newWebhookRequest(ctx context.Context, method, uri string, body map[string]string) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, uri))
	if err != nil {
		return nil, err
	}

	var buf io.Reader
	if body != nil {
		form := url.Values{}
		form.Set("client_id", c.stravaKey)
		form.Set("client_secret", c.stravaSecret)
		for key, value := range body {
			form.Set(key, value)
		}
		buf = ioutil.NopCloser(bytes.NewBufferString(form.Encode()))
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
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
