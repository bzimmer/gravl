package cyclinganalytics

//go:generate go run ../../cmd/genwith/genwith.go --do --auth --package cyclinganalytics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl/pkg"
)

const (
	baseURL = "https://www.cyclinganalytics.com/api"
)

// Client .
type Client struct {
	config oauth2.Config
	token  oauth2.Token
	client *http.Client

	User  *UserService
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
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	c.User = &UserService{client: c}
	c.Rides = &RidesService{client: c}

	return c, nil
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
	req.Header.Set("User-Agent", pkg.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token.AccessToken))
	return req, nil
}
