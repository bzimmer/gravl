package zwift

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bzimmer/gravl/pkg/providers/activity"
	"golang.org/x/oauth2"
)

//go:generate genwith --do --client --token --package zwift

const baseURL = "https://us-or-rly101.zwift.com"
const userAgent = "CNL/3.4.1 (Darwin Kernel 20.3.0) zwift/1.0.61590 curl/7.64.1"

// Endpoint is Zwifts's OAuth 2.0 endpoint
var Endpoint = oauth2.Endpoint{
	TokenURL:  "https://secure.zwift.com/auth/realms/zwift/tokens/access/codes",
	AuthStyle: oauth2.AuthStyleAutoDetect,
}

// Client for communicating with Zwift
type Client struct {
	token  *oauth2.Token
	client *http.Client

	Auth     *AuthService
	Activity *ActivityService
	Profile  *ProfileService
}

func (c *Client) Exporter() activity.Exporter {
	return c.Activity
}

func withServices() Option {
	return func(c *Client) error {
		c.Auth = &AuthService{c}
		c.Profile = &ProfileService{c}
		c.Activity = &ActivityService{c}
		c.token.TokenType = "bearer"
		return nil
	}
}

func (c *Client) newAPIRequest(ctx context.Context, method, uri string) (*http.Request, error) {
	if c.token.AccessToken == "" {
		return nil, errors.New("accessToken required")
	}
	q := fmt.Sprintf("%s/%s", baseURL, uri)
	u, err := url.Parse(q)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token.AccessToken))
	return req, nil
}
