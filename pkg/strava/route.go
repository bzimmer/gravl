package strava

import (
	"context"
	"fmt"
	"net/http"

	"github.com/twpayne/go-polyline"

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
	rtes := make([]*Route, count)
	uri := fmt.Sprintf("athletes/%d/routes?page=%d&per_page=%d", p.athleteID, start, count)
	req, err := p.service.client.newAPIRequest(p.ctx, http.MethodGet, uri)
	if err != nil {
		return 0, err
	}
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
//  call with (ctx, total, start, count)
func (s *RouteService) Routes(ctx context.Context, athleteID int, specs ...int) ([]*route.Route, error) {
	p := &routePaginator{
		service:   *s,
		athleteID: athleteID,
		ctx:       ctx,
		routes:    make([]*route.Route, 0),
	}
	err := paginate(p, specs...)
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
	c, _, err := polyline.DecodeCoords([]byte(r.Map.Polyline))
	if err != nil {
		return nil, err
	}
	zero := float64(0)
	coords := make([][]float64, len(c))
	for i, x := range c {
		coords[i] = []float64{x[1], x[0], zero}
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
