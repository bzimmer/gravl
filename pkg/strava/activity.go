package strava

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ActivityService .
type ActivityService service

const (
	pageSize = 100
)

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

// Streams of data from the activity
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
