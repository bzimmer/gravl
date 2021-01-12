package cyclinganalytics

import (
	"context"
	"net/http"
)

// UserService .
type UserService service

func (s *UserService) Me(ctx context.Context) (*User, error) {
	uri := "me"
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}

	usr := &User{}
	err = s.client.do(req, usr)
	if err != nil {
		return nil, err
	}
	return usr, nil
}
