package cyclinganalytics

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

// RidesService manages rides for a user
type RidesService service

// RideOptions specify additional detail to return for a queried ride
type RideOptions struct {
	// Streams is a list of valid data streams
	Streams []string
	// Curves specifies the difference curves to return
	Curves struct {
		// AveragePower will be returned if true
		AveragePower bool
		// EffectivePower will be returned if true
		EffectivePower bool
	}
}

const (
	meupload = "me/upload"
)

func (r *RideOptions) values() (*url.Values, error) {
	v := &url.Values{}
	if r.Streams != nil {
		if err := validateStreams(r.Streams); err != nil {
			return nil, err
		}
		v.Set("streams", strings.Join(r.Streams, ","))
	}
	if r.Curves.AveragePower && r.Curves.EffectivePower {
		v.Set("curves", "true")
	} else {
		v.Set("power_curve", fmt.Sprintf("%t", r.Curves.AveragePower))
		v.Set("epower_curve", fmt.Sprintf("%t", r.Curves.EffectivePower))
	}
	return v, nil
}

// Ride returns a single ride with available options
func (s *RidesService) Ride(ctx context.Context, rideID int64, opts RideOptions) (*Ride, error) {
	uri := fmt.Sprintf("ride/%d", rideID)
	params, err := opts.values()
	if err != nil {
		return nil, err
	}
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, params, nil)
	if err != nil {
		return nil, err
	}
	ride := &Ride{}
	err = s.client.do(req, ride)
	if err != nil {
		return nil, err
	}
	return ride, nil
}

// Rides returns a slice of rides for the user
func (s *RidesService) Rides(ctx context.Context, userID UserID, spec activity.Pagination) ([]*Ride, error) {
	uri := "me/rides"
	if userID != Me {
		uri = fmt.Sprintf("%d/rides", userID)
	}
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	type rides struct {
		Rides []*Ride `json:"rides"`
	}
	res := &rides{}
	err = s.client.do(req, res)
	if err != nil {
		return nil, err
	}
	if spec.Total > 0 {
		n := math.Min(float64(len(res.Rides)), float64(spec.Total))
		res.Rides = res.Rides[:int(n)]
	}
	return res.Rides, nil
}

// Upload the file for the authenticated user
func (s *RidesService) Upload(ctx context.Context, file *activity.File) (*Upload, error) {
	return s.UploadWithUser(ctx, Me, file)
}

// Upload the file for the user
func (s *RidesService) UploadWithUser(ctx context.Context, userID UserID, file *activity.File) (*Upload, error) {
	if file == nil {
		return nil, errors.New("missing upload file")
	}

	uri := meupload
	if userID != Me {
		uri = fmt.Sprintf("user/%d/upload", userID)
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if err := w.WriteField("filename", file.Name); err != nil {
		return nil, err
	}
	fw, err := w.CreateFormFile("data", file.Name)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}

	req, err := s.client.newAPIRequest(ctx, http.MethodPost, uri, nil, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	res := &Upload{}
	err = s.client.do(req, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Status returns the status of an upload request
func (s *RidesService) Status(ctx context.Context, uploadID int64) (*Upload, error) {
	return s.StatusWithUser(ctx, Me, uploadID)
}

// Status returns the status of an upload request for the user
func (s *RidesService) StatusWithUser(ctx context.Context, userID UserID, uploadID int64) (*Upload, error) {
	uri := meupload
	if userID != Me {
		uri = fmt.Sprintf("user/%d/upload", userID)
	}
	uri = fmt.Sprintf("%s/%d", uri, uploadID)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	res := &Upload{}
	err = s.client.do(req, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AvailableStreams returns the list of valid stream names
func (s *RidesService) StreamSets() map[string]string {
	q := make(map[string]string)
	for k, v := range streamsets {
		q[k] = v
	}
	return q
}

func validateStreams(streams []string) error {
	for i := range streams {
		_, ok := streamsets[streams[i]]
		if !ok {
			return fmt.Errorf("invalid stream '%s'", streams[i])
		}
	}
	return nil
}

// https://www.cyclinganalytics.com/developer/api#/ride/ride_id
var streamsets = map[string]string{
	"cadence":                "",
	"distance":               "The sequence of distance values for this stream, in kilometers [float]",
	"elevation":              "The sequence of elevation values for this stream, in meters [float]",
	"gears":                  "",
	"gradient":               "The sequence of grade values for this stream, as percents of a grade [float]",
	"heart_rate_variability": "",
	"heartrate":              "The sequence of heart rate values for this stream, in beats per minute [integer]",
	"latitude":               "",
	"longitude":              "",
	"lrbalance":              "",
	"pedal_smoothness":       "",
	"platform_center_offset": "",
	"power_direction":        "",
	"power_phase":            "",
	"power":                  "",
	"respiration_rate":       "",
	"smo2":                   "",
	"speed":                  "The sequence of speed values for this stream, in meters per second [float]",
	"temperature":            "The sequence of temperature values for this stream, in celsius degrees [float]",
	"thb":                    "",
	"torque_effectiveness":   "",
}
