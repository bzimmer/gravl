package strava

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

// RouteService is the API for route endpoints
type RouteService service

type routePaginator struct {
	athleteID int
	routes    []*Route
	service   RouteService
}

func (p *routePaginator) PageSize() int {
	return PageSize
}

func (p *routePaginator) Count() int {
	return len(p.routes)
}

func (p *routePaginator) Do(ctx context.Context, spec activity.Pagination) (int, error) {
	uri := fmt.Sprintf("athletes/%d/routes?page=%d&per_page=%d", p.athleteID, spec.Start, spec.Count)
	req, err := p.service.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return 0, err
	}
	var rts []*Route
	err = p.service.client.do(req, &rts)
	if err != nil {
		return 0, err
	}
	if len(p.routes)+len(rts) > spec.Total {
		rts = rts[:spec.Total-len(p.routes)]
	}
	p.routes = append(p.routes, rts...)
	return len(rts), nil
}

// Routes returns a page of routes for an athlete
func (s *RouteService) Routes(ctx context.Context, athleteID int, spec activity.Pagination) ([]*Route, error) {
	p := &routePaginator{service: *s, athleteID: athleteID, routes: make([]*Route, 0)}
	err := activity.Paginate(ctx, p, spec)
	if err != nil {
		return nil, err
	}
	return p.routes, nil
}

// Route returns a route
func (s *RouteService) Route(ctx context.Context, routeID int64) (*Route, error) {
	uri := fmt.Sprintf("routes/%d", routeID)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	rte := &Route{}
	err = s.client.do(req, &rte)
	if err != nil {
		return nil, err
	}
	return rte, nil
}
