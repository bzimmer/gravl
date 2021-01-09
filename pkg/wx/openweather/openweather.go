package openweather

//go:generate go run ../../../cmd/genwith/genwith.go --do --auth --client --package openweather

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
	baseURL = "http://api.openweathermap.org/data/2.5"
)

// Client .
type Client struct {
	config oauth2.Config
	token  oauth2.Token
	client *http.Client

	Forecast *ForecastService
}

func withServices() Option {
	return func(c *Client) error {
		c.Forecast = &ForecastService{client: c}
		return nil
	}
}

func (c *Client) newAPIRequest(ctx context.Context, method, uri string, values *url.Values) (*http.Request, error) {
	if c.token.AccessToken == "" {
		return nil, errors.New("accessToken required")
	}
	values.Set("mode", "json")
	values.Set("appid", c.token.AccessToken)
	u, err := url.Parse(fmt.Sprintf("%s/%s?%s", baseURL, uri, values.Encode()))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", pkg.UserAgent)
	return req, nil
}
