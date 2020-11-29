package rwgps

//go:generate go run ../../cmd/genwith/genwith.go --do --auth --package rwgps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl/pkg"
)

const (
	apiVersion = "2"
	baseURL    = "https://ridewithgps.com"
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
	req.Header.Set("User-Agent", pkg.UserAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
