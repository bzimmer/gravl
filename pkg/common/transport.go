package common

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// VerboseTransport .
type VerboseTransport struct {
	Event     *zerolog.Event
	Transport http.RoundTripper
}

// RoundTrip .
func (t *VerboseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := http.DefaultTransport
	if t.Transport != nil {
		transport = t.Transport
	}
	event := log.Debug()
	if t.Event != nil {
		event = t.Event
	}
	dump, _ := httputil.DumpRequestOut(req, true)
	event.Str("req", string(dump)).Msg("sending")
	res, err := transport.RoundTrip(req)
	dump, _ = httputil.DumpResponse(res, true)
	event.Str("res", string(dump)).Msg("received")
	return res, err
}

// TestDataTransport .
type TestDataTransport struct {
	Status      int
	Filename    string
	ContentType string
	RequestFunc func(*http.Request) error
}

// RoundTrip .
func (t *TestDataTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		err  error
		data []byte
	)
	if t.RequestFunc != nil {
		err = t.RequestFunc(req)
		if err != nil {
			return nil, err
		}
	}
	if t.Filename != "" {
		filename := filepath.Join("testdata", t.Filename)
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	} else {
		data = make([]byte, 0)
	}

	header := make(http.Header)
	header.Add("Content-Type", t.ContentType)

	return &http.Response{
		StatusCode:    t.Status,
		ContentLength: int64(len(data)),
		Body:          ioutil.NopCloser(bytes.NewBuffer(data)),
		Header:        header,
		Request:       req,
	}, nil
}
