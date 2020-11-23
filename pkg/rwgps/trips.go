package rwgps

import (
	"context"
	"fmt"
	"net/http"
)

// TripsService .
type TripsService service

// Trip .
func (s *TripsService) Trip(ctx context.Context, tripID int64) (*Trip, error) {
	return s.trip(ctx, OriginTrip, fmt.Sprintf("trips/%d.json", tripID))
}

// Route .
func (s *TripsService) Route(ctx context.Context, routeID int64) (*Trip, error) {
	return s.trip(ctx, OriginRoute, fmt.Sprintf("routes/%d.json", routeID))
}

func (s *TripsService) trip(ctx context.Context, origin Origin, uri string) (*Trip, error) {
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	res := &TripResponse{}
	err = s.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	var t *Trip
	switch origin {
	case OriginTrip:
		t = res.Trip
	case OriginRoute:
		t = res.Route
	default:
		return nil, fmt.Errorf("unknown origin type {%d}", origin)
	}
	t.Origin = origin
	return t, nil
}
