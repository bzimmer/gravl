package rwgps

//go:generate go run ../../dev/genwith.go --auth --package rwgps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const (
	apiVersion = "2"
	baseURL    = "https://ridewithgps.com"
	userAgent  = "(github.com/bzimmer/gravl/rwgps)"
)

// https://ridewithgps.com/api?lang=en

// Client .
type Client struct {
	config oauth2.Config
	token  oauth2.Token
	client *http.Client

	Users *UsersService
	Trips *TripsService
}

type service struct {
	client *Client //nolint:golint,structcheck
}

// Option .
type Option func(*Client) error

var headers = map[string]string{
	"User-Agent":   userAgent,
	"Accept":       "application/json",
	"Content-type": "application/json"}

// NewClient .
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		client: &http.Client{},
		token:  oauth2.Token{},
		config: oauth2.Config{},
	}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	c.Users = &UsersService{client: c}
	c.Trips = &TripsService{client: c}

	return c, nil
}

func (c *Client) newAPIRequest(ctx context.Context, method, uri string) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, uri))
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(
		map[string]string{
			"version":    apiVersion,
			"apikey":     c.config.ClientID,
			"auth_token": c.token.AccessToken,
		})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
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
