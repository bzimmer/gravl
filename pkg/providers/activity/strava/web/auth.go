package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gocolly/colly/v2"
)

// AuthService is the API for auth endpoints
type AuthService service

// Login creates a user session for the username/password
func (s *AuthService) Login(ctx context.Context, username, password string) error {
	if s.client.client.Jar == nil {
		return errors.New("cookiejar not set on http client")
	}
	s.client.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() {
		// cookies are not set properly if redirects are enabled
		// login on the main thread, establish the proper cookies, then use this client in more than one thread
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
	if param == "" || token == "" {
		return errors.New("one of param or token is nil")
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
