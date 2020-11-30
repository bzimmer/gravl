package cyclinganalytics

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// RidesService .
type RidesService service

type RideOptions struct {
	Streams []string
	Curves  struct {
		AveragePower   bool
		EffectivePower bool
	}
}

func (r *RideOptions) values() *url.Values {
	v := &url.Values{}
	if r.Streams != nil {
		v.Set("streams", strings.Join(r.Streams, ","))
	}
	if r.Curves.AveragePower && r.Curves.EffectivePower {
		v.Set("curves", "true")
	} else {
		v.Set("power_curve", fmt.Sprintf("%t", r.Curves.AveragePower))
		v.Set("epower_curve", fmt.Sprintf("%t", r.Curves.EffectivePower))
	}
	return v
}

func (s *RidesService) Ride(ctx context.Context, rideID int64, opts RideOptions) (*Ride, error) {
	uri := fmt.Sprintf("ride/%d", rideID)

	params := opts.values()
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, params)
	if err != nil {
		return nil, err
	}
	ride := &Ride{}
	err = s.client.do(req, ride)
	if err != nil {
		return nil, err
	}
	return ride, nil
}

func (s *RidesService) Rides(ctx context.Context, userID UserID) ([]*Ride, error) {
	uri := "me/rides"
	if userID == Me {
		uri = fmt.Sprintf("%d/rides", userID)
	}
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Rides []*Ride `json:"rides"`
	}
	res := &Response{}
	err = s.client.do(req, res)
	if err != nil {
		return nil, err
	}
	return res.Rides, nil
}
