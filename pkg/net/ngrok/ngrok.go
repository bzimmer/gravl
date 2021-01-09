package ngrok

//go:generate go run ../../../cmd/genwith/genwith.go --client --do --package ngrok

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const baseURL = "http://127.0.0.1:4040/api"

// Client for accessing ngrok endpoints
type Client struct {
	client *http.Client

	Tunnels *TunnelsService
}

func withServices() Option {
	return func(c *Client) error {
		c.Tunnels = &TunnelsService{c}
		return nil
	}
}

func (c *Client) newRequest(ctx context.Context, method, uri string) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, uri))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
