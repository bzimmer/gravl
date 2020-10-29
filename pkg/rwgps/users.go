package rwgps

import (
	"context"
	"net/http"
)

// UsersService .
type UsersService service

// AuthenticatedUser .
func (s *UsersService) AuthenticatedUser(ctx context.Context) (*User, error) {
	uri := "users/current.json"
	req, err := s.client.newAPIRequest(http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	res := &UserResponse{}
	err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return res.User, err
}