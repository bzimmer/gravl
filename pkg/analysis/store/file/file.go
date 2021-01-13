package file

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

type file struct {
	path       string
	activities map[int64]*strava.Activity
}

// Open a file of json-encoded activities
func Open(path string) (store.Source, error) {
	return &file{path: path}, nil
}

// Close the file
func (s *file) Close() error {
	return nil
}

// Activities returns a slice of (potentially incomplete) Activity instances
func (s *file) Activities(ctx context.Context) ([]*strava.Activity, error) {
	var b []byte
	var err error
	var sc fastjson.Scanner

	s.activities = make(map[int64]*strava.Activity)

	b, err = ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	sc.InitBytes(b)
	for sc.Next() {
		if err = sc.Error(); err != nil {
			return nil, err
		}
		val := sc.Value()
		act := &strava.Activity{}
		err = json.Unmarshal(val.MarshalTo(nil), act)
		if err != nil {
			return nil, err
		}
		s.activities[act.ID] = act
	}

	var i int
	acts := make([]*strava.Activity, len(s.activities))
	for _, act := range s.activities {
		acts[i] = act
		i++
	}
	return acts, nil
}

// Activity returns a fully populated Activity
func (s *file) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	act, ok := s.activities[activityID]
	if !ok {
		return nil, fmt.Errorf("id not found {%d}", activityID)
	}
	return act, nil
}
