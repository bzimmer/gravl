package strava

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ActivityService is the API for activity endpoints
type ActivityService service

type activityPaginator struct {
	service    ActivityService
	activities []*Activity
}

func (p *activityPaginator) count() int {
	return len(p.activities)
}

func (p *activityPaginator) do(ctx context.Context, start, count int) (int, error) {
	uri := fmt.Sprintf("athlete/activities?page=%d&per_page=%d", start, count)
	req, err := p.service.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return 0, err
	}
	var acts []*Activity
	err = p.service.client.do(req, &acts)
	if err != nil {
		return 0, err
	}
	for _, act := range acts {
		if err != nil {
			return 0, err
		}
		p.activities = append(p.activities, act)
	}
	return len(acts), nil
}

// Streams returns the activities data streams
func (s *ActivityService) Streams(ctx context.Context, activityID int64, streams ...string) (*Streams, error) {
	keys := strings.Join(streams, ",")
	uri := fmt.Sprintf("activities/%d/streams/%s?key_by_type=true", activityID, keys)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	sts := make(map[string]*Stream)
	err = s.client.do(req, &sts)
	if err != nil {
		return nil, err
	}
	return &Streams{ActivityID: activityID, Streams: sts}, err
}

// Activity returns the activity specified by id for an athlete
func (s *ActivityService) Activity(ctx context.Context, id int64) (*Activity, error) {
	uri := fmt.Sprintf("activities/%d", id)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	act := &Activity{}
	err = s.client.do(req, act)
	if err != nil {
		return nil, err
	}
	return act, err
}

// Activities returns a page of activities for an athlete
func (s *ActivityService) Activities(ctx context.Context, spec Pagination) ([]*Activity, error) {
	p := &activityPaginator{
		service:    *s,
		activities: make([]*Activity, 0),
	}
	err := paginate(ctx, p, spec)
	if err != nil {
		return nil, err
	}
	return p.activities, nil
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
