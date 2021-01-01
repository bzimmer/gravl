package strava

//go:generate go run ../../cmd/genwith/genwith.go --do --client --endpoint --auth --ratelimit --package strava

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl/pkg"
)

const (
	baseURL    = "https://www.strava.com/api/v3"
	baseWebURL = "https://www.strava.com/"
	PageSize   = 100
)

// Endpoint is Strava's OAuth 2.0 endpoint
var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://www.strava.com/oauth/authorize",
	TokenURL:  "https://www.strava.com/oauth/token",
	AuthStyle: oauth2.AuthStyleAutoDetect,
}

// Client client
type Client struct {
	client *http.Client
	config oauth2.Config
	token  oauth2.Token

	Auth     *AuthService
	Route    *RouteService
	Fitness  *FitnessService
	Webhook  *WebhookService
	Athlete  *AthleteService
	Activity *ActivityService
}

func withServices(c *Client) {
	c.Auth = &AuthService{client: c}
	c.Route = &RouteService{client: c}
	c.Fitness = &FitnessService{client: c}
	c.Webhook = &WebhookService{client: c}
	c.Athlete = &AthleteService{client: c}
	c.Activity = &ActivityService{client: c}
}

func WithCookieJar() Option {
	return func(c *Client) error {
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			return err
		}
		c.client.Jar = jar
		return nil
	}
}

func (c *Client) newWebRequest(ctx context.Context, method, uri string, values url.Values) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", baseWebURL, uri))
	if err != nil {
		return nil, err
	}
	var b io.Reader
	if values != nil {
		b = strings.NewReader(values.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), b)
	if err != nil {
		return nil, err
	}
	// req.Header.Set("User-Agent", pkg.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	if values != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, nil
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
	u, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, uri))
	if err != nil {
		return nil, err
	}

	var buf io.Reader
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
