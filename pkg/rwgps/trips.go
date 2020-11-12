package rwgps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bzimmer/gravl/pkg/common/route"
)

// TripsService .
type TripsService service

const (
	tripType  = "trip"
	routeType = "route"
)

// Trip .
func (s *TripsService) Trip(ctx context.Context, tripID int64) (*route.Route, error) {
	return s.trip(ctx, tripType, fmt.Sprintf("trips/%d.json", tripID))
}

// Route .
func (s *TripsService) Route(ctx context.Context, routeID int64) (*route.Route, error) {
	return s.trip(ctx, routeType, fmt.Sprintf("routes/%d.json", routeID))
}

func (s *TripsService) trip(ctx context.Context, activity, uri string) (*route.Route, error) {
	req, err := s.client.newAPIRequest(http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	res := &TripResponse{}
	err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}

	var t *Trip
	switch activity {
	case tripType:
		t = res.Trip
	case routeType:
		t = res.Route
	default:
		return nil, fmt.Errorf("unknown activity type {%s}", activity)
	}
	return newRoute(activity, t)
}

func newRoute(activity string, trip *Trip) (*route.Route, error) {
	coords := make([][]float64, len(trip.TrackPoints))
	for i, tp := range trip.TrackPoints {
		coords[i] = []float64{tp.Longitude, tp.Latitude, tp.Elevation}
	}
	return &route.Route{
		ID:          fmt.Sprintf("%d", trip.ID),
		Name:        trip.Name,
		Source:      baseURL,
		Origin:      routeOrigin(activity),
		Description: trip.Description,
		Coordinates: coords,
	}, nil
}

func routeOrigin(activity string) route.Origin {
	switch activity {
	case tripType:
		return route.Activity
	case routeType:
		return route.Planned
	}
	return route.Unknown
}
