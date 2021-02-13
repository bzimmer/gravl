package cyclinganalytics

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

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
	polls           = 5
	pollingDuration = 2 * time.Second
	meupload        = "me/upload"
)

func (r *RideOptions) values() *url.Values {
	v := &url.Values{}
	if r.Streams != nil {
		v.Set("streams", strings.Join(r.Streams, ","))
	}
	if r.Curves.AveragePower && r.Curves.EffectivePower {
		v.Set("curves", "true")
	} else {
		v.Set("power_curve", fmt.Sprintf("%t", r.Curves.AveragePower))
		v.Set("epower_curve", fmt.Sprintf("%t", r.Curves.EffectivePower))
	}
	return v
}

// Ride returns a single ride with available options
func (s *RidesService) Ride(ctx context.Context, rideID int64, opts RideOptions) (*Ride, error) {
	uri := fmt.Sprintf("ride/%d", rideID)

	params := opts.values()
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
	res := &RidesResponse{}
	err = s.client.do(req, res)
	if err != nil {
		return nil, err
	}
	return res.Rides, nil
}

// File for uploading
type File struct {
	io.Reader
	Name string
}

func (f *File) Close() error {
	if f.Reader == nil {
		return nil
	}
	if x, ok := f.Reader.(io.Closer); ok {
		return x.Close()
	}
	return nil
}

// Upload the file for the user
func (s *RidesService) Upload(ctx context.Context, userID UserID, file *File) (*Upload, error) {
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
func (s *RidesService) Status(ctx context.Context, userID UserID, uploadID int64) (*Upload, error) {
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

// Poll the status of an upload
//
// The operation will continue until either it is completed (status != "processing"), the context
//  is canceled, or the maximum number of iterations have been exceeded.
//
// More information can be found at:
//  https://www.cyclinganalytics.com/developer/api#/user/user_id/upload/upload_id
func (s *RidesService) Poll(ctx context.Context, userID UserID, uploadID int64) <-chan *UploadResult {
	res := make(chan *UploadResult)
	go func() {
		defer close(res)
		i := 0
		for ; i < polls; i++ {
			upload, err := s.Status(ctx, userID, uploadID)
			if err != nil {
				res <- &UploadResult{Err: err}
				return
			}
			res <- &UploadResult{Upload: upload}
			// status: processing, done, or error
			if upload.Status != "processing" {
				return
			}
			select {
			case <-ctx.Done():
				res <- &UploadResult{Err: ctx.Err()}
				return
			case <-time.After(pollingDuration):
			}
		}
		if i == polls {
			log.Warn().Int("polls", polls).Msg("exceeded max polls")
		}
	}()
	return res
}

// // ValidStream returns true if the strean name is valid
// func ValidStream(stream string) bool { // nolint
// 	// https://www.cyclinganalytics.com/developer/api#/ride/ride_id
// 	switch stream {
// 	case "cadence":
// 		return true
// 	case "distance":
// 		// The sequence of distance values for this stream, in kilometers [float]
// 		return true
// 	case "elevation":
// 		// The sequence of elevation values for this stream, in meters [float]
// 		return true
// 	case "gears":
// 		return true
// 	case "gradient":
// 		// The sequence of grade values for this stream, as percents of a grade [float]
// 		return true
// 	case "heart_rate_variability":
// 		return true
// 	case "heartrate":
// 		// The sequence of heart rate values for this stream, in beats per minute [integer]
// 		return true
// 	case "latitude":
// 		return true
// 	case "longitude":
// 		return true
// 	case "lrbalance":
// 		return true
// 	case "pedal_smoothness":
// 		return true
// 	case "platform_center_offset":
// 		return true
// 	case "power":
// 		return true
// 	case "power_direction":
// 		return true
// 	case "power_phase":
// 		return true
// 	case "respiration_rate":
// 		return true
// 	case "smo2":
// 		return true
// 	case "speed":
// 		// The sequence of speed values for this stream, in meters per second [float]
// 		return true
// 	case "temperature":
// 		// The sequence of temperature values for this stream, in celsius degrees [float]
// 		return true
// 	case "thb":
// 		return true
// 	case "torque_effectiveness":
// 		return true
// 	}
// 	return false
// }
