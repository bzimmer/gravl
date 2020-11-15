package strava

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/twpayne/go-polyline"

	"github.com/bzimmer/gravl/pkg/common/route"
)

// ActivityService .
type ActivityService service

type activityPaginator struct {
	service    ActivityService
	ctx        context.Context
	activities []*route.Route
}

// Count .
func (p *activityPaginator) Count() int {
	return len(p.activities)
}

// Do .
func (p *activityPaginator) Do(start, count int) (int, error) {
	uri := fmt.Sprintf("athlete/activities?page=%d&per_page=%d", start, count)
	req, err := p.service.client.newAPIRequest(http.MethodGet, uri)
	if err != nil {
		return 0, err
	}
	acts := make([]*Activity, count)
	err = p.service.client.Do(p.ctx, req, &acts)
	if err != nil {
		return 0, err
	}
	for _, act := range acts {
		r, err := newRouteFromActivity(act)
		if err != nil {
			return 0, err
		}
		p.activities = append(p.activities, r)
	}
	return len(acts), nil
}

// Route returns a Route from an activities stream data
//  This is different from returning a Strava Route
func (s *ActivityService) Route(ctx context.Context, activityID int64) (*route.Route, error) {
	streams, err := s.Streams(ctx, activityID, "latlng", "elevation")
	if err != nil {
		return nil, err
	}
	return newRouteFromStreams(activityID, streams)
}

// Streams returns the activities data streams
func (s *ActivityService) Streams(ctx context.Context, activityID int64, streams ...string) (map[string]*Stream, error) {
	keys := strings.Join(streams, ",")
	uri := fmt.Sprintf("activities/%d/streams/%s?key_by_type=true", activityID, keys)
	req, err := s.client.newAPIRequest(http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*Stream)
	err = s.client.Do(ctx, req, &m)
	if err != nil {
		return nil, err
	}
	return m, err
}

// Activity returns the activity specified by id for an athlete
func (s *ActivityService) Activity(ctx context.Context, id int64) (*Activity, error) {
	uri := fmt.Sprintf("activities/%d", id)
	req, err := s.client.newAPIRequest(http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	act := &Activity{}
	err = s.client.Do(ctx, req, act)
	if err != nil {
		return nil, err
	}
	return act, err
}

// Activities returns a page of activities for an athlete
//  call with (ctx, total, start, count)
func (s *ActivityService) Activities(ctx context.Context, specs ...int) ([]*route.Route, error) {
	p := &activityPaginator{
		service:    *s,
		ctx:        ctx,
		activities: make([]*route.Route, 0),
	}
	err := paginate(p, specs...)
	if err != nil {
		return nil, err
	}
	return p.activities, nil
}

func newRouteFromActivity(a *Activity) (*route.Route, error) {
	zero := float64(0)
	var coords [][]float64
	// unclear which of the polylines will be valid
	for _, p := range []string{a.Map.Polyline, a.Map.SummaryPolyline} {
		if p != "" {
			c, _, err := polyline.DecodeCoords([]byte(p))
			if err != nil {
				return nil, err
			}
			coords = make([][]float64, len(c))
			for i, x := range c {
				coords[i] = []float64{x[1], x[0], zero}
			}
			break
		}
	}
	rte := &route.Route{
		ID:          fmt.Sprintf("%d", a.ID),
		Name:        a.Name,
		Description: a.Description,
		Source:      baseURL,
		Origin:      route.Activity,
		Coordinates: coords,
	}
	return rte, nil
}

func newRouteFromStreams(activityID int64, streams map[string]*Stream) (*route.Route, error) {
	latlng, ok := streams["latlng"]
	if !ok {
		return nil, errors.New("missing required latlng stream")
	}

	zero := float64(0)
	rte := &route.Route{
		ID:          fmt.Sprintf("%d", activityID),
		Source:      baseURL,
		Origin:      route.Activity,
		Coordinates: make([][]float64, len(latlng.Data)),
	}

	altitude, ok := streams["altitude"]
	for i, m := range latlng.Data {
		lat := m.([]interface{})[0]
		lng := m.([]interface{})[1]
		alt := zero
		if ok {
			alt = (altitude.Data[i]).(float64)
		}
		rte.Coordinates[i] = []float64{lng.(float64), lat.(float64), alt}
	}
	return rte, nil
}

// // https://developers.strava.com/docs/reference/#api-models-StreamSet
// func validStream(name string) bool {
// 	switch name {
// 	case "latlng":
// 		// The sequence of lat/long values for this stream [float, float]
// 		return true
// 	case "altitude":
// 		// The sequence of altitude values for this stream, in meters [float]
// 		return true
// 	case "time":
// 		// The sequence of time values for this stream, in seconds [integer]
// 		return true
// 	case "distance":
// 		// The sequence of distance values for this stream, in meters [float]
// 		return true
// 	case "velocity_smooth":
// 		// The sequence of velocity values for this stream, in meters per second [float]
// 		return true
// 	case "heartrate":
// 		// The sequence of heart rate values for this stream, in beats per minute [integer]
// 		return true
// 	case "cadence":
// 		// The sequence of cadence values for this stream, in rotations per minute [integer]
// 		return true
// 	case "watts":
// 		// The sequence of power values for this stream, in watts [integer]
// 		return true
// 	case "temp":
// 		// The sequence of temperature values for this stream, in celsius degrees [float]
// 		return true
// 	case "moving":
// 		// The sequence of moving values for this stream, as boolean values [boolean]
// 		return true
// 	case "grade_smooth":
// 		// The sequence of grade values for this stream, as percents of a grade [float]
// 		return true
// 	}
// 	return false
// }
