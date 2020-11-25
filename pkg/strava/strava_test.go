package strava_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bzimmer/gravl/pkg/strava"
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
		strava.WithHTTPTracing(true),
		strava.WithTokenCredentials("fooKey", "barToken", time.Time{}))
}

type ManyTransport struct {
	Filename string
}

func (t *ManyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	q := req.URL.Query()
	n, _ := strconv.Atoi(q.Get("per_page"))

	data, err := ioutil.ReadFile(t.Filename)
	if err != nil {
		return nil, err
	}

	acts := make([]string, 0)
	for i := 0; i < n; i++ {
		acts = append(acts, string(data))
	}

	res := strings.Join(acts, ",")
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString("[" + res + "]")),
		Header:     make(http.Header),
	}, nil
}
