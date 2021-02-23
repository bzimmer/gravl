package strava

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bzimmer/gravl/pkg/providers/activity"
)

// ActivityService is the API for activity endpoints
type ActivityService service

const (
	polls           = 5
	pollingDuration = 2 * time.Second
)

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
	req, err := p.service.client.newAPIRequest(ctx, http.MethodGet, uri, nil)
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
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil)
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
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil)
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

// Upload the file for the user
//
// More information can be found at https://developers.strava.com/docs/uploads/
func (s *ActivityService) Upload(ctx context.Context, file *activity.File) (*Upload, error) {
	if file == nil || file.Name == "" || file.Format == activity.Original {
		return nil, errors.New("missing upload file, name, or format")
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if err := w.WriteField("filename", file.Name); err != nil {
		return nil, err
	}
	if err := w.WriteField("data_type", file.Format.String()); err != nil {
		return nil, err
	}
	fw, err := w.CreateFormFile("file", file.Name)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}

	req, err := s.client.newAPIRequest(ctx, http.MethodPost, "uploads", &b)
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
//
// More information can be found at https://developers.strava.com/docs/uploads/
func (s *ActivityService) Status(ctx context.Context, uploadID int64) (*Upload, error) {
	uri := fmt.Sprintf("uploads/%d", uploadID)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri, nil)
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
// The operation will continue until either it is completed, the context
//  is canceled, or the maximum number of iterations have been exceeded.
//
// More information can be found at:
//   https://developers.strava.com/docs/uploads/
//   A successful upload will return a response with an upload ID. You may use this ID to poll the
//   status of your upload. Strava recommends polling no more than once a second. The mean processing
//   time is around 8 seconds.
func (s *ActivityService) Poll(ctx context.Context, uploadID int64) <-chan *UploadResult {
	res := make(chan *UploadResult)
	go func() {
		defer close(res)
		i := 0
		for ; i < polls; i++ {
			var r *UploadResult
			upload, err := s.Status(ctx, uploadID)
			switch {
			case err != nil:
				r = &UploadResult{Err: err}
			default:
				r = &UploadResult{Upload: upload}
			}
			select {
			case <-ctx.Done():
				log.Error().Err(ctx.Err()).Msg("ctx is done")
				return
			case res <- r:
				if upload.ActivityID > 0 || upload.Error != "" {
					return
				}
			case <-time.After(pollingDuration):
			}
		}
		if i == polls {
			log.Warn().Int("polls", polls).Msg("exceeded max polls")
		}
	}()
	return res
}

// // ValidStream returns true if the stream name is valid
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
