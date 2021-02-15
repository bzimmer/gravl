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

type channelPaginator struct {
	service    ActivityService
	count      int
	activities chan *ActivityResult
}

func (p *channelPaginator) PageSize() int {
	return PageSize
}

func (p *channelPaginator) Count() int {
	return p.count
}

func (p *channelPaginator) Do(ctx context.Context, spec activity.Pagination) (int, error) {
	uri := fmt.Sprintf("athlete/activities?page=%d&per_page=%d", spec.Start, spec.Count)
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
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case p.activities <- &ActivityResult{Activity: act}:
			p.count++
		}
		if p.count == spec.Total {
			break
		}
	}
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

// ActivityResult is the result of querying for a stream of activities
type ActivityResult struct {
	Activity *Activity
	Err      error
}

// Activities returns a channel for activities and errors for an athlete
//
// Either the first error or last activity will close the channel
func (s *ActivityService) Activities(ctx context.Context, spec activity.Pagination) <-chan *ActivityResult {
	acts := make(chan *ActivityResult, PageSize)
	go func() {
		defer close(acts)
		p := &channelPaginator{service: *s, activities: acts}
		err := activity.Paginate(ctx, p, spec)
		if err != nil {
			acts <- &ActivityResult{Err: err}
		}
	}()
	return acts
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
