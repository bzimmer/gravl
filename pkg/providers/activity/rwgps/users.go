package rwgps

import (
	"context"
)

// UsersService provides access to the user API
type UsersService service

// AuthenticatedUser returns the authenticated user
func (s *UsersService) AuthenticatedUser(ctx context.Context) (*User, error) {
	uri := "users/current.json"
	req, err := s.client.newAPIRequest(ctx, uri, nil)
	if err != nil {
		return nil, err
	}
	type UserResponse struct {
		User *User `json:"user"`
	}
	res := &UserResponse{}
	err = s.client.do(req, res)
	if err != nil {
		return nil, err
	}
	return res.User, err
}
