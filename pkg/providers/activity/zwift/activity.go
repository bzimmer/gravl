package zwift

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/bzimmer/gravl/pkg/providers/activity"
	"github.com/rs/zerolog/log"
)

// ActivityService is the API for profile endpoints
type ActivityService service

// pageSize default for querying bulk entities (eg trips, routes)
const pageSize = 20

type paginator struct {
	service    ActivityService
	athleteID  int64
	activities []*Activity
}

func (p *paginator) PageSize() int {
	return pageSize
}

func (p *paginator) Count() int {
	return len(p.activities)
}

func (p *paginator) Do(ctx context.Context, spec activity.Pagination) (int, error) {
	// pagination uses the concept of page (based on strava), rwgps uses an offset by row
	//  since pagination starts with page 1 (again, strava), subtract one from `start`
	count := spec.Count
	start := int64((spec.Start - 1) * p.PageSize())

	uri := fmt.Sprintf("/api/profiles/%d/activities/?start=%d&limit=%d", p.athleteID, start, count)
	req, err := p.service.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return 0, err
	}
	var acts []*Activity
	if err = p.service.client.do(req, &acts); err != nil {
		return 0, err
	}
	if spec.Total > 0 && len(p.activities)+len(acts) > spec.Total {
		acts = acts[:spec.Total-len(p.activities)]
	}
	p.activities = append(p.activities, acts...)
	return len(acts), nil
}

func (s *ActivityService) Activity(ctx context.Context, athleteID int64, activityID int64) (*Activity, error) {
	uri := fmt.Sprintf("/api/profiles/%d/activities/%d", athleteID, activityID)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	var act *Activity
	if err = s.client.do(req, &act); err != nil {
		if x, ok := err.(*Fault); ok {
			if x.Message == "" {
				x.Message = "Zwift does not use JSON errors, enable http tracing for more details"
			}
		}
		return nil, err
	}
	return act, nil
}

func (s *ActivityService) Activities(ctx context.Context, athleteID int64, spec activity.Pagination) ([]*Activity, error) {
	p := &paginator{service: *s, athleteID: athleteID, activities: make([]*Activity, 0)}
	err := activity.Paginate(ctx, p, spec)
	if err != nil {
		return nil, err
	}
	return p.activities, nil
}

func (s *ActivityService) Export(ctx context.Context, act *Activity) (*activity.Export, error) {
	uri := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", act.FitFileBucket, act.FitFileKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	res, err := s.client.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}
	defer res.Body.Close()
	if res.StatusCode >= http.StatusBadRequest {
		if res.StatusCode == http.StatusNotFound {
			return nil, &Fault{Message: "activity not found"}
		}
		return nil, &Fault{Message: fmt.Sprintf("error code: %d", res.StatusCode)}
	}
	out := &bytes.Buffer{}
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return nil, err
	}
	/*
		The parse fails with "expected slash after first token" without the "attachment"
		prefix added. The HTTP headers below are captured from a query.

		HTTP/1.1 200 OK
		Content-Length: 101003
		Accept-Ranges: bytes
		Content-Disposition: filename=2020-11-24-07-28-35.fit
		Content-Type: application/octet-stream
		Date: Thu, 18 Feb 2021 00:38:11 GMT
		Etag: "4665f697cd0685029cc3ac4fc08d35ef"
		Last-Modified: Tue, 24 Nov 2020 16:14:06 GMT
		Server: AmazonS3
	*/
	disposition := res.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType("attachment; " + disposition)
	if err != nil {
		return nil, err
	}
	log.Info().Str("filename", params["filename"]).Int64("activityID", act.ID).Msg("export")
	return &activity.Export{
		Reader: out,
		ID:     act.ID,
		Name:   params["filename"],
		Format: activity.FIT,
	}, nil
}
