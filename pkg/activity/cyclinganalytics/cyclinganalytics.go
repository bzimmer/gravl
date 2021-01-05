package cyclinganalytics

//go:generate go run ../../cmd/genwith/genwith.go --do --client --endpoint --auth --package cyclinganalytics

import (
	"context"
	"errors"
	"fmt"
	"io"
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

// Endpoint is CyclingAnalytics's OAuth 2.0 endpoint
var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://www.cyclinganalytics.com/api/auth",
	TokenURL:  "https://www.cyclinganalytics.com/api/token",
	AuthStyle: oauth2.AuthStyleAutoDetect,
}

func withServices(c *Client) error {
	c.User = &UserService{client: c}
	c.Rides = &RidesService{client: c}
	return nil
}

func (c *Client) newAPIRequest(ctx context.Context, method, uri string, values *url.Values, body io.Reader) (*http.Request, error) {
	if c.token.AccessToken == "" {
		return nil, errors.New("accessToken required")
	}
	q := fmt.Sprintf("%s/%s", baseURL, uri)
	if values != nil {
		q = fmt.Sprintf("%s?%s", q, values.Encode())
	}
	u, err := url.Parse(q)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", pkg.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token.AccessToken))
	return req, nil
}
