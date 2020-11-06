package strava

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	gj "github.com/paulmach/go.geojson"
	"github.com/rs/zerolog/log"
)

// ActivityService .
type ActivityService service

const (
	pageSize = 100
)

// Streams of data from the activity
func (s *ActivityService) Streams(ctx context.Context, activityID int64, streams ...string) (*gj.FeatureCollection, error) {
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
	return newFeatureCollection(activityID, m)
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
func (s *ActivityService) Activities(ctx context.Context, specs ...int) (*[]Activity, error) {
	var start, count, total int
	switch len(specs) {
	case 0:
		total, start, count = 0, 1, pageSize
	case 1:
		total, start, count = specs[0], 1, pageSize
	case 2:
		total, start, count = specs[0], specs[1], pageSize
	case 3:
		total, start, count = specs[0], specs[1], specs[2]
	default:
		return nil, errors.New("too many varargs")
	}
	if total < 0 {
		return nil, errors.New("total less than zero")
	}
	if total <= count {
		count = total
	}
	return s.activities(ctx, total, start, count)
}

func (s *ActivityService) activities(ctx context.Context, total, start, count int) (*[]Activity, error) {
	all := make([]Activity, 0)

	for {
		acts := make([]Activity, count)
		uri := fmt.Sprintf("athlete/activities?page=%d&per_page=%d", start, count)
		req, err := s.client.newAPIRequest(http.MethodGet, uri)
		if err != nil {
			return nil, err
		}
		err = s.client.Do(ctx, req, &acts)
		if err != nil {
			return nil, err
		}
		for _, act := range acts {
			all = append(all, act)
		}
		if len(acts) != count || len(all) >= total {
			break
		}
		start = start + 1
		if (total - len(all)) < pageSize {
			count = total - len(all)
		} else {
			count = pageSize
		}
	}

	return &all, nil
}

func newFeatureCollection(activityID int64, streams map[string]*Stream) (*gj.FeatureCollection, error) {
	fc := gj.NewFeatureCollection()

	if streams == nil {
		return fc, nil
	}

	// The sequence of lat/long values for this stream
	latlng, ok := streams["latlng"]
	if !ok {
		return nil, errors.New("missing latlng stream")
	}
	delete(streams, "latlng")

	n := len(latlng.Data)
	log.Debug().Str("name", "latlng").Int("count", n).Msg("fc")
	for name, stream := range streams {
		if n != len(stream.Data) {
			return nil, errors.New("inconsistent streams sizes")
		}
		log.Debug().Str("name", name).Int("count", len(stream.Data)).Msg("fc")
	}

	coords := make([][]float64, n)
	feature := gj.NewFeature(gj.NewLineStringGeometry(coords))
	feature.ID = activityID

	zero := float64(0)
	// The sequence of altitude values for this stream, in meters
	altitude, ok := streams["altitude"]
	for i, m := range latlng.Data {
		lat := m.([]interface{})[0]
		lng := m.([]interface{})[1]
		alt := zero
		if ok {
			alt = (altitude.Data[i]).(float64)
		}
		coords[i] = []float64{lng.(float64), lat.(float64), alt}
	}
	delete(streams, "altitude")

	dataStreams := make(map[string]interface{})
	for name, stream := range streams {
		n, s := dataStream(name, stream)
		dataStreams[n] = s
	}
	feature.Properties["streams"] = dataStreams

	fc.AddFeature(feature)
	return fc, nil
}

func dataStream(name string, stream *Stream) (string, []interface{}) {
	switch name {
	case "time":
		// The sequence of time values for this stream, in seconds [integer]
	case "distance":
		// The sequence of distance values for this stream, in meters [float]
	case "velocity_smooth":
		// The sequence of velocity values for this stream, in meters per second [float]
	case "heartrate":
		// The sequence of heart rate values for this stream, in beats per minute [integer]
	case "cadence":
		// The sequence of cadence values for this stream, in rotations per minute [integer]
	case "watts":
		// The sequence of power values for this stream, in watts [integer]
	case "temp":
		// The sequence of temperature values for this stream, in celsius degrees [float]
	case "moving":
		// The sequence of moving values for this stream, as boolean values [boolean]
	case "grade_smooth":
		// The sequence of grade values for this stream, as percents of a grade [float]
	}
	return name, stream.Data
}
