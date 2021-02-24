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
	if err := validateStreams(streams); err != nil {
		return nil, err
	}
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
	if len(streams) > 0 {
		// confirm valid streams before querying strava for the activity
		if err := validateStreams(streams); err != nil {
			return nil, err
		}
	}
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
	// @todo(bzimmer) query streams concurrently to activity
	if len(streams) > 0 {
		var sms *Streams
		log.Debug().Strs("streams", streams).Msg("querying")
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
			log.Error().Err(err).Msg("paginating activities")
			select {
			case <-ctx.Done():
				return
			case acts <- &ActivityResult{Err: err}:
				return
			}
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

// https://developers.strava.com/docs/reference/#api-models-StreamSet
var streamsets = map[string]string{
	"altitude":        "The sequence of altitude values for this stream, in meters [float]",
	"cadence":         "The sequence of cadence values for this stream, in rotations per minute [integer]",
	"distance":        "The sequence of distance values for this stream, in meters [float]",
	"grade_smooth":    "The sequence of grade values for this stream, as percents of a grade [float]",
	"heartrate":       "The sequence of heart rate values for this stream, in beats per minute [integer]",
	"latlng":          "The sequence of lat/long values for this stream [float, float]",
	"moving":          "The sequence of moving values for this stream, as boolean values [boolean]",
	"temp":            "The sequence of temperature values for this stream, in celsius degrees [float]",
	"time":            "The sequence of time values for this stream, in seconds [integer]",
	"velocity_smooth": "The sequence of velocity values for this stream, in meters per second [float]",
	"watts":           "The sequence of power values for this stream, in watts [integer]",
}

// AvailableStreams returns the list of valid stream names
func (s *ActivityService) StreamSets() map[string]string {
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
