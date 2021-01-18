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
func (s *file) Activities(ctx context.Context) (<-chan *strava.Activity, <-chan error) {
	errs := make(chan error)
	acts := make(chan *strava.Activity)

	s.activities = make(map[int64]*strava.Activity)

	go func() {
		defer close(acts)
		defer close(errs)

		var b []byte
		var err error
		var sc fastjson.Scanner
		b, err = ioutil.ReadFile(s.path)
		if err != nil {
			errs <- err
			return
		}
		sc.InitBytes(b)
		for sc.Next() {
			if err = sc.Error(); err != nil {
				errs <- err
				return
			}
			val := sc.Value()
			act := &strava.Activity{}
			err = json.Unmarshal(val.MarshalTo(nil), act)
			if err != nil {
				errs <- err
				return
			}
			s.activities[act.ID] = act
			acts <- act
		}
	}()
	return acts, errs
}

// Activity returns a fully populated Activity
func (s *file) Activity(ctx context.Context, activityID int64) (*strava.Activity, error) {
	act, ok := s.activities[activityID]
	if !ok {
		return nil, fmt.Errorf("id not found {%d}", activityID)
	}
	return act, nil
}
