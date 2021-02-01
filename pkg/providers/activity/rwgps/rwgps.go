package rwgps

//go:generate genwith --do --client --token --config --package rwgps

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
	// PageSize default for querying bulk entities (eg trips, routes)
	PageSize = 30
)

// Client for communicating with RWGPS
type Client struct {
	config oauth2.Config
	token  *oauth2.Token
	client *http.Client

	Users *UsersService
	Trips *TripsService
}

func withServices() Option {
	return func(c *Client) error {
		c.Users = &UsersService{client: c}
		c.Trips = &TripsService{client: c}
		return nil
	}
}

func (c *Client) newAPIRequest(ctx context.Context, method, uri string, params map[string]string) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, uri))
	if err != nil {
		return nil, err
	}
	x := map[string]string{
		"version":    apiVersion,
		"apikey":     c.config.ClientID,
		"auth_token": c.token.AccessToken,
	}
	for k, v := range params {
		x[k] = v
	}
	b, err := json.Marshal(x)
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
