package visualcrossing

//go:generate go run ../../cmd/genwith/genwith.go --auth --client --package visualcrossing

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl/pkg"
)

const (
	baseURL = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/weatherdata"
)

// Client .
type Client struct {
	config oauth2.Config
	token  oauth2.Token
	client *http.Client

	Forecast *ForecastService
}

func withServices(c *Client) {
	c.Forecast = &ForecastService{client: c}
}

func (c *Client) newAPIRequest(ctx context.Context, method, uri string, values *url.Values) (*http.Request, error) {
	if c.token.AccessToken == "" {
		return nil, errors.New("accessToken required")
	}
	v := url.Values{
		"key":              []string{c.token.AccessToken},
		"contentType":      []string{"json"},
		"locationMode":     []string{"array"},
		"shortColumnNames": []string{"false"},
	}
	for key, val := range *values {
		v[key] = val
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s?%s", baseURL, uri, v.Encode()))
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

// do executes the request
func (c *Client) do(req *http.Request, v interface{}) error {
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

	var (
		buf    bytes.Buffer
		fault  Fault
		reader = io.TeeReader(res.Body, &buf)
	)

	// VC uses StatusOK for everything, sigh
	err = json.NewDecoder(reader).Decode(&fault)
	if err != nil {
		return err
	} else if code := fault.ErrorCode; code != 0 && code != http.StatusOK {
		return &fault
	}

	if v != nil {
		err := json.NewDecoder(&buf).Decode(v)
		if err != nil {
			return err
		}
	}
	return nil
}
