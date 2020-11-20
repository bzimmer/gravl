package strava

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bzimmer/gravl/pkg/common/route"
)

// RouteService is the API for route endpoints
type RouteService service

type routePaginator struct {
	athleteID int
	routes    []*route.Route
	ctx       context.Context
	service   RouteService
}

// Count .
func (p *routePaginator) Count() int {
	return len(p.routes)
}

// Do .
func (p *routePaginator) Do(start, count int) (int, error) {
	uri := fmt.Sprintf("athletes/%d/routes?page=%d&per_page=%d", p.athleteID, start, count)
	req, err := p.service.client.newAPIRequest(p.ctx, http.MethodGet, uri)
	if err != nil {
		return 0, err
	}
	var rtes []*Route
	err = p.service.client.Do(req, &rtes)
	if err != nil {
		return 0, err
	}
	for _, rte := range rtes {
		r, err := newRouteFromRoute(rte)
		if err != nil {
			return 0, err
		}
		p.routes = append(p.routes, r)
	}
	return len(rtes), nil
}

// Routes returns a page of routes for an athlete
func (s *RouteService) Routes(ctx context.Context, athleteID int, spec Pagination) ([]*route.Route, error) {
	p := &routePaginator{
		service:   *s,
		athleteID: athleteID,
		ctx:       ctx,
		routes:    make([]*route.Route, 0),
	}
	err := paginate(p, spec)
	if err != nil {
		return nil, err
	}
	return p.routes, nil
}

// Route .
func (s *RouteService) Route(ctx context.Context, routeID int64) (*route.Route, error) {
	uri := fmt.Sprintf("routes/%d", routeID)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	rte := &Route{}
	err = s.client.Do(req, &rte)
	if err != nil {
		return nil, err
	}
	return newRouteFromRoute(rte)
}

func newRouteFromRoute(r *Route) (*route.Route, error) {
	coords, err := polylineToCoords(r.Map.Polyline, r.Map.SummaryPolyline)
	if err != nil {
		return nil, err
	}
	rte := &route.Route{
		ID:          fmt.Sprintf("%d", r.ID),
		Name:        r.Name,
		Description: r.Description,
		Source:      baseURL,
		Origin:      route.Planned,
		Coordinates: coords,
	}
	return rte, nil
}
