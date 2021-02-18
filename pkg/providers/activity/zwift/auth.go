package zwift

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// AuthService is the API for auth endpoints
type AuthService service

// Zwift uses `expires_in` rather than `expires`
type expiresin struct {
	oauth2.Token
	Expiry int `json:"expires_in,omitempty"`
}

func (s *AuthService) Refresh(ctx context.Context, username, password string) (*oauth2.Token, error) {
	values := url.Values{
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
		"client_id":  {"Zwift_Mobile_Link"},
	}
	body := strings.NewReader(values.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, Endpoint.TokenURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var token *expiresin
	if err = s.client.do(req, &token); err != nil {
		return nil, err
	}
	token.Token.Expiry = time.Now().Add(time.Duration(token.Expiry) * time.Second)
	return &token.Token, nil
}
