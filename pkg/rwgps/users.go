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
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	res := &UserResponse{}
	err = s.client.do(req, res)
	if err != nil {
		return nil, err
	}
	return res.User, err
}
