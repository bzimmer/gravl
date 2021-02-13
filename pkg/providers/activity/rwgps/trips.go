package rwgps

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

// pageSize default for querying bulk entities (eg trips, routes)
const pageSize = 100

// TripsService provides access to Trips and Routes via the RWGPS API
type TripsService service

type tripsPaginator struct {
	service TripsService
	userID  UserID
	trips   []*Trip
}

func (p *tripsPaginator) Page() int {
	return pageSize
}

func (p *tripsPaginator) Count() int {
	return len(p.trips)
}

func (p *tripsPaginator) Do(ctx context.Context, start, count int) (int, error) {
	uri := fmt.Sprintf("users/%d/trips.json", p.userID)
	params := map[string]string{
		// pagination uses the concept of page (based on strava), rwgps uses an offset by row
		//  since pagination starts with page 1 (again, strava), subtract one from `start`
		"offset": strconv.FormatInt(int64((start-1)*p.Page()), 10),
		"limit":  strconv.FormatInt(int64(count), 10),
	}
	req, err := p.service.client.newAPIRequest(ctx, http.MethodGet, uri, params)
	if err != nil {
		return 0, err
	}
	type TripsResponse struct {
		Results      []*Trip `json:"results"`
		ResultsCount int     `json:"results_count"`
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
func (s *TripsService) Trips(ctx context.Context, userID UserID, spec activity.Pagination) ([]*Trip, error) {
	p := &tripsPaginator{service: *s, userID: userID, trips: make([]*Trip, 0)}
	err := activity.Paginate(ctx, p, spec)
	if err != nil {
		return nil, err
	}
	return p.trips, nil
}

// Trip returns a trip for the `tripID`
func (s *TripsService) Trip(ctx context.Context, tripID int64) (*Trip, error) {
	return s.trip(ctx, TypeTrip, fmt.Sprintf("trips/%d.json", tripID))
}

// Route returns a trip for the `routeID`
func (s *TripsService) Route(ctx context.Context, routeID int64) (*Trip, error) {
	return s.trip(ctx, TypeRoute, fmt.Sprintf("routes/%d.json", routeID))
}

func (s *TripsService) trip(ctx context.Context, entity Type, uri string) (*Trip, error) {
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	type TripResponse struct {
		Type  string `json:"type"`
		Trip  *Trip  `json:"trip"`
		Route *Trip  `json:"route"`
	}

	res := &TripResponse{}
	err = s.client.do(req, res)
	if err != nil {
		return nil, err
	}

	var t *Trip
	switch entity {
	case TypeTrip:
		t = res.Trip
	case TypeRoute:
		t = res.Route
	}
	t.Type = entity.String()
	return t, nil
}
