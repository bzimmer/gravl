package rwgps

import (
	"context"
	"fmt"
	"net/http"
)

// TripsService provides access to Trips and Routes via the RWGPS API
type TripsService service

type tripsPaginator struct {
	service TripsService
	userID  UserID
	trips   []*Trip
}

func (p *tripsPaginator) page() int {
	return PageSize
}

func (p *tripsPaginator) count() int {
	return len(p.trips)
}

func (p *tripsPaginator) do(ctx context.Context, start, count int) (int, error) {
	uri := fmt.Sprintf("users/%d/trips.json", p.userID)
	params := map[string]string{
		"offset": fmt.Sprintf("%d", start),
		"limit":  fmt.Sprintf("%d", count),
	}
	req, err := p.service.client.newAPIRequest(ctx, http.MethodGet, uri, params)
	if err != nil {
		return 0, err
	}
	res := &TripsResponse{}
	err = p.service.client.do(req, res)
	if err != nil {
		return 0, err
	}
	p.trips = append(p.trips, res.Results...)
	return len(res.Results), nil
}

// Trips returns a slice of Trips
func (s *TripsService) Trips(ctx context.Context, userID UserID, spec Pagination) ([]*Trip, error) {
	p := &tripsPaginator{
		service: *s,
		userID:  userID,
		trips:   make([]*Trip, 0),
	}
	err := paginate(ctx, p, spec)
	if err != nil {
		return nil, err
	}
	return p.trips, nil
}

// Trip returns a trip for the `tripID`
func (s *TripsService) Trip(ctx context.Context, tripID int64) (*Trip, error) {
	return s.trip(ctx, OriginTrip, fmt.Sprintf("trips/%d.json", tripID))
}

// Route returns a trip for the `routeID`
func (s *TripsService) Route(ctx context.Context, routeID int64) (*Trip, error) {
	return s.trip(ctx, OriginRoute, fmt.Sprintf("routes/%d.json", routeID))
}

func (s *TripsService) trip(ctx context.Context, origin Origin, uri string) (*Trip, error) {
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	res := &TripResponse{}
	err = s.client.do(req, res)
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
