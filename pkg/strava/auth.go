package strava

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gocolly/colly/v2"
	"golang.org/x/oauth2"
)

// AuthService is the API for auth endpoints
type AuthService service

// Refresh returns a new access token
func (s *AuthService) Refresh(ctx context.Context) (*oauth2.Token, error) {
	t := s.client.config.TokenSource(ctx, &s.client.token)
	t = oauth2.ReuseTokenSource(&s.client.token, t)
	return t.Token()
}

func (s *AuthService) Login(ctx context.Context, username, password string) error {
	// inspired by https://github.com/pR0Ps/stravaweblib
	if s.client.client.Jar == nil {
		return errors.New("cookiejar not set on http client")
	}
	s.client.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() {
		// @todo(bzimmer) this is a hack -- the function is set on the client which is shared among threads
		s.client.client.CheckRedirect = nil
	}()

	var param, token string
	c := colly.NewCollector(colly.AllowedDomains("www.strava.com"))
	c.SetClient(s.client.client)
	c.OnHTML("meta", func(e *colly.HTMLElement) {
		switch e.Attr("name") {
		case "csrf-param":
			param = e.Attr("content")
		case "csrf-token":
			token = e.Attr("content")
		}
	})
	if err := c.Visit(fmt.Sprintf("%s/login", baseWebURL)); err != nil {
		return err
	}

	v := url.Values{}
	v.Set("remember_me", "on")
	v.Set("email", username)
	v.Set("password", password)
	v.Set(param, token)

	req, err := s.client.newWebRequest(ctx, http.MethodPost, "session", v)
	if err != nil {
		return err
	}
	res, err := s.client.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
	}
	defer res.Body.Close()
	return nil
}
