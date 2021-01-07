package strava

//go:generate go run ../../../cmd/genwith/genwith.go --do --client --endpoint --auth --ratelimit --package strava

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl/pkg"
)

const (
	baseURL = "https://www.strava.com/api/v3"
	// PageSize default for querying bulk entities (eg activities, routes)
	PageSize = 100
)

// Endpoint is Strava's OAuth 2.0 endpoint
var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://www.strava.com/oauth/authorize",
	TokenURL:  "https://www.strava.com/oauth/token",
	AuthStyle: oauth2.AuthStyleAutoDetect,
}

// Client for communicating with Strava
type Client struct {
	client *http.Client
	config oauth2.Config
	token  oauth2.Token

	Auth     *AuthService
	Route    *RouteService
	Webhook  *WebhookService
	Athlete  *AthleteService
	Activity *ActivityService
}

func withServices(c *Client) error { // nolint:unparam
	c.Auth = &AuthService{client: c}
	c.Route = &RouteService{client: c}
	c.Webhook = &WebhookService{client: c}
	c.Athlete = &AthleteService{client: c}
	c.Activity = &ActivityService{client: c}
	return nil
}

func (c *Client) newAPIRequest(ctx context.Context, method, uri string) (*http.Request, error) { // nolint:unparam
	if c.token.AccessToken == "" {
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
	req.Header.Set("User-Agent", pkg.UserAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.token.AccessToken))
	return req, nil
}

func (c *Client) newWebhookRequest(ctx context.Context, method, uri string, body map[string]string) (*http.Request, error) {
	var buf io.Reader
	u, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, uri))
	if err != nil {
		return nil, err
	}
	if body != nil {
		form := url.Values{}
		form.Set("client_id", c.config.ClientID)
		form.Set("client_secret", c.config.ClientSecret)
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
