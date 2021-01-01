// The web package scrapes the strava website for functionality not available via the API.
//
// Inspired by https://github.com/pR0Ps/stravaweblib
package web

//go:generate go run ../../../cmd/genwith/genwith.go --client --package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

const baseWebURL = "https://www.strava.com/"

// Client client
type Client struct {
	client *http.Client

	Auth    *AuthService
	Fitness *FitnessService
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

func withServices(c *Client) error {
	if c.client.Jar == nil {
		return errors.New("no cookiejar set; use WithHTTPClient() or WithCookieJar()")
	}
	c.Auth = &AuthService{client: c}
	c.Fitness = &FitnessService{client: c}
	return nil
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

// do executes the http request and populates v with the result.
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
	if res.StatusCode >= http.StatusBadRequest {
		return errors.New("something bad happened; enable --http-tracing to investigate")
	}
	if v != nil {
		err := json.NewDecoder(res.Body).Decode(v)
		if err == io.EOF {
			err = nil // ignore EOF errors caused by empty response body
		}
		return err
	}
	return nil
}
