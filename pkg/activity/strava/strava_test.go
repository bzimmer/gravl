package strava_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bzimmer/gravl/pkg/activity/strava"
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
	n, _ := strconv.Atoi(q.Get("per_page"))
	if t.Total > 0 {
		n = t.Total
		if t.total >= t.Total {
			n = 0
		}
	}

	data, err := ioutil.ReadFile(t.Filename)
	if err != nil {
		return nil, err
	}

	var acts []string
	for i := 0; i < n; i++ {
		acts = append(acts, string(data))
	}
	t.total = t.total + len(acts)

	res := strings.Join(acts, ",")
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString("[" + res + "]")),
		Header:     make(http.Header),
	}, nil
}
