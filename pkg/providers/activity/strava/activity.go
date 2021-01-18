package strava

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

// ActivityService is the API for activity endpoints
type ActivityService service

var _ activity.Paginator = &slicePaginator{}
var _ activity.Paginator = &channelPaginator{}

type channelPaginator struct {
	service    ActivityService
	count      int
	activities chan *Activity
}

func (p *channelPaginator) Page() int {
	return PageSize
}

func (p *channelPaginator) Count() int {
	return p.count
}

func (p *channelPaginator) Do(ctx context.Context, start, count int) (int, error) {
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
		p.count++
		p.activities <- act
	}
	return len(acts), nil
}

type slicePaginator struct {
	service    ActivityService
	activities []*Activity
}

func (p *slicePaginator) Page() int {
	return PageSize
}

func (p *slicePaginator) Count() int {
	return len(p.activities)
}

func (p *slicePaginator) Do(ctx context.Context, start, count int) (int, error) {
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
	p.activities = append(p.activities, acts...)
	return len(acts), nil
}

// Streams returns the activity's data streams
func (s *ActivityService) Streams(ctx context.Context, activityID int64, streams ...string) (*Streams, error) {
	keys := strings.Join(streams, ",")
	uri := fmt.Sprintf("activities/%d/streams/%s?key_by_type=true", activityID, keys)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	sts := &Streams{}
	err = s.client.do(req, sts)
	if err != nil {
		return nil, err
	}
	sts.ActivityID = activityID
	return sts, err
}

// Activity returns the activity specified by id
func (s *ActivityService) Activity(ctx context.Context, activityID int64, streams ...string) (*Activity, error) {
	uri := fmt.Sprintf("activities/%d", activityID)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	act := &Activity{}
	err = s.client.do(req, act)
	if err != nil {
		return nil, err
	}
	if len(streams) > 0 {
		var sms *Streams
		log.Debug().Strs("streams", streams).Msg("querying")
		// @todo(bzimmer) query streams concurrently to activity
		sms, err = s.Streams(ctx, activityID, streams...)
		if err != nil {
			return nil, err
		}
		act.Streams = sms
	}
	return act, err
}

// Activities returns channels for activities and errors for an athlete
//
// Either the first error or last activity will close the channels
func (s *ActivityService) Activities(ctx context.Context, spec activity.Pagination) (<-chan *Activity, <-chan error) {
	errs := make(chan error)
	acts := make(chan *Activity)
	go func() {
		defer close(acts)
		defer close(errs)
		p := &channelPaginator{service: *s, activities: acts}
		err := activity.Paginate(ctx, p, spec)
		if err != nil {
			errs <- err
		}
	}()
	return acts, errs
}

// Activities returns a slice of activities from the channels returned by `ActivityService.Activities`
//
// This is a convenience function for those operations which require the full activity slice.
func Activities(ctx context.Context, acts <-chan *Activity, errs <-chan error) ([]*Activity, error) {
	var activities []*Activity
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err, ok := <-errs:
			// if the channel was not closed an error occurred so return it
			// if the channel is closed do nothing to ensure the activity channel can run to
			//  completion and return the full slice of activities
			if ok {
				return nil, err
			}
		case act, ok := <-acts:
			if !ok {
				// the channel is closed, return the activities
				return activities, nil
			}
			activities = append(activities, act)
		}
	}
}

// // ValidStream returns true if the strean name is valid
// func ValidStream(stream string) bool {
// 	// https://developers.strava.com/docs/reference/#api-models-StreamSet
// 	switch stream {
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
