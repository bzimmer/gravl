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

func (r *RideOptions) Values() (*url.Values, error) {
	v := &url.Values{}
	if r.Streams != nil {
		v.Set("streams", strings.Join(r.Streams, ","))
	}
	if r.Curves.AveragePower && r.Curves.EffectivePower {
		v.Set("curves", "true")
	} else {
		if r.Curves.AveragePower {
			v.Set("power_curve", "true")
		}
		if r.Curves.EffectivePower {
			v.Set("epower_curve", "true")
		}
	}
	return v, nil
}

func (s *RidesService) Ride(ctx context.Context, rideID int64, opts RideOptions) (*Ride, error) {
	uri := fmt.Sprintf("ride/%d", rideID)

	params, err := opts.Values()
	if err != nil {
		return nil, err
	}
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, params)
	if err != nil {
		return nil, err
	}
	ride := &Ride{}
	err = s.client.Do(req, ride)
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
	err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}
	return res.Rides, nil
}
