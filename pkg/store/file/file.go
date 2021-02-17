package file

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/bzimmer/gravl/pkg/store"
	"github.com/bzimmer/gravl/pkg/store/memory"
)

type Option func(*file) error

type file struct {
	path  string
	flush bool
}

func (f *file) Activities() ([]*strava.Activity, error) {
	var b []byte
	var err error
	var activities []*strava.Activity
	defer func(t time.Time) {
		log.Info().
			Dur("elapsed", time.Since(t)).
			Str("path", f.path).
			Int("activities", len(activities)).
			Msg("file")
	}(time.Now())
	b, err = ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	err = nil
	gjson.ForEachLine(string(b), func(res gjson.Result) bool {
		act := &strava.Activity{}
		err = json.Unmarshal([]byte(res.Raw), act)
		if err != nil {
			return false
		}
		activities = append(activities, act)
		return true
	})
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (f *file) Close(activities map[int64]*strava.Activity) error {
	if f.flush {
		fo, err := ioutil.TempFile("", "")
		if err != nil {
			return err
		}
		defer os.Remove(fo.Name())
		writer := bufio.NewWriter(fo)
		enc := json.NewEncoder(writer)
		enc.SetIndent("", " ")
		enc.SetEscapeHTML(false)
		for i := range activities {
			if err := enc.Encode(activities[i]); err != nil {
				return err
			}
		}
		if err := writer.Flush(); err != nil {
			return err
		}
		if err := fo.Close(); err != nil {
			return err
		}
		if err := os.Rename(fo.Name(), f.path); err != nil {
			return err
		}
		return nil
	}
	return nil
}

// Flush configures whether to flush the contents of any updates (add or remove) to disk on close
func Flush(flush bool) Option {
	return func(f *file) error {
		f.flush = flush
		return nil
	}
}

// Open a file of json-encoded activities
func Open(path string, opts ...Option) (store.Store, error) {
	f := &file{path: path}
	for i := range opts {
		if err := opts[i](f); err != nil {
			return nil, err
		}
	}
	return memory.Open(f)
}
