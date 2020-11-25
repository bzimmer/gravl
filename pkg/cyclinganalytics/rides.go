package cyclinganalytics

import (
	"context"
	"net/http"
)

// RidesService .
type RidesService service

func (s *RidesService) Rides(ctx context.Context) ([]*Ride, error) {
	uri := "me/rides"
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Rides []*Ride `json:"rides"`
	}
	res := &Response{}
	err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}
	return res.Rides, nil
}
