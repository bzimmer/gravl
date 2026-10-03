package cyclinganalytics_test

import (
	"encoding/json"
	"net/http"
	"testing"

	api "github.com/bzimmer/activity/cyclinganalytics"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/internal"
)

func TestAPIFailures(t *testing.T) {
	for _, args := range [][]string{
		{"gravl", "cyclinganalytics", "athlete"},
		{"gravl", "cyclinganalytics", "activities"},
		{"gravl", "cyclinganalytics", "activity", "77282721"},
	} {
		t.Run(args[2], func(t *testing.T) {
			tt := &internal.Harness{Name: args[2], Args: args, AnyErr: true}
			internal.Run(t, tt, internal.FailingAPI(t), command)
		})
	}
}

func TestEncodingFailures(t *testing.T) {
	a := assert.New(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/me", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(&api.User{}))
	})
	mux.HandleFunc("/me/rides", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(struct {
			Rides []*api.Ride `json:"rides"`
		}{Rides: []*api.Ride{{ID: 1008}}}))
	})
	mux.HandleFunc("/ride/77282721", func(w http.ResponseWriter, _ *http.Request) {
		a.NoError(json.NewEncoder(w).Encode(&api.Ride{ID: 77282721}))
	})
	for _, args := range [][]string{
		{"gravl", "cyclinganalytics", "athlete"},
		{"gravl", "cyclinganalytics", "activities"},
		{"gravl", "cyclinganalytics", "activity", "77282721"},
		{"gravl", "cyclinganalytics", "streamsets"},
	} {
		t.Run(args[2], func(t *testing.T) {
			tt := &internal.Harness{Name: args[2], Args: args, Before: internal.FailEncoding, Err: internal.ErrEncode.Error()}
			internal.Run(t, tt, mux, command)
		})
	}
}

func TestInvalidRideID(t *testing.T) {
	tt := &internal.Harness{
		Name: "invalid id",
		Args: []string{"gravl", "cyclinganalytics", "activity", "not-a-number"},
		Err:  "invalid syntax",
	}
	internal.Run(t, tt, internal.FailingAPI(t), command)
}
