package strava

import (
	"context"

	"golang.org/x/oauth2"
)

// AuthService is the API for auth endpoints
type AuthService service

// Refresh returns a new access token
func (s *AuthService) Refresh(ctx context.Context) (*oauth2.Token, error) {
	t := s.client.config.TokenSource(ctx, s.client.token)
	t = oauth2.ReuseTokenSource(s.client.token, t)
	return t.Token()
}
