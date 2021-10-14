package strava_test

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
	"github.com/bzimmer/httpwares"
)

func newClient(status int, filename string) (*strava.Client, error) {
	return newClienter(status, filename, nil, nil)
}

func newClienter(status int, filename string, requester httpwares.Requester, responder httpwares.Responder) (*strava.Client, error) {
	return strava.NewClient(
		strava.WithTransport(&httpwares.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "application/json",
			Requester:   requester,
			Responder:   responder}),
		strava.WithHTTPTracing(false),
		strava.WithTokenCredentials("fooKey", "barToken", time.Time{}))
}

type ManyTransport struct {
	Filename string
	Total    int
	total    int
}

func (t *ManyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	q := req.URL.Query()
	n, err := strconv.Atoi(q.Get("per_page"))
	if err != nil {
		return nil, err
	}
	if t.Total > 0 {
		n = t.Total
		if t.total >= t.Total {
			n = 0
		}
	}
	t.total = t.total + n

	data, err := os.ReadFile(t.Filename)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	if err := body.WriteByte(byte('[')); err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		if _, err := body.Write(data); err != nil {
			return nil, err
		}
		if i+1 < n {
			if err := body.WriteByte(','); err != nil {
				return nil, err
			}
		}
	}
	if err := body.WriteByte(byte(']')); err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(&body),
		Header:     make(http.Header),
	}, nil
}
