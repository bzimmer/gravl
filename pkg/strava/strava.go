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
	"time"

	"github.com/bzimmer/gravl/pkg/common"
	"github.com/markbates/goth"
)

const (
	stravaURI = "https://www.strava.com/api/v3"
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
	Webhook  *WebhookService
	Athlete  *AthleteService
	Activity *ActivityService
}

type service struct {
	client *Client
}

// Option .
type Option func(*Client) error

// WithHTTPTracing .
func WithHTTPTracing(debug bool) Option {
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
	return func(client *Client) error {
		client.client.Timeout = timeout
		return nil
	}
}

// WithWebhookCredentials provides the Strava credentials
func WithWebhookCredentials(stravaKey, stravaSecret string) Option {
	return func(client *Client) error {
		client.stravaKey = stravaKey
		client.stravaSecret = stravaSecret
		return nil
	}
}

// WithAPICredentials provides the Strava credentials
func WithAPICredentials(accessToken, refreshToken string) Option {
	return func(client *Client) error {
		client.accessToken = accessToken
		client.refreshToken = refreshToken
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
	c.Webhook = &WebhookService{client: c}
	c.Athlete = &AthleteService{client: c}
	c.Activity = &ActivityService{client: c}

	return c, nil
}

func (c *Client) newAPIRequest(method, uri string) (*http.Request, error) {
	if c.accessToken == "" {
		return nil, errors.New("accessToken required")
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s", stravaURI, uri))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", c.accessToken))
	return req, nil
}

func (c *Client) newWebhookRequest(method, uri string, body map[string]string) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", stravaURI, uri))
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

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
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
