package common

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestDataTransport struct {
	status   int
	filename string
}

func (r *TestDataTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var data []byte
	if r.filename != "" {
		dir, _ := os.Getwd()
		filename := filepath.Join(dir, "../../testdata", r.filename)
		data, _ = ioutil.ReadFile(filename)
	} else {
		data = make([]byte, 0)
	}

	header := make(http.Header)
	header.Add("content-type", "text/html; charset=utf-8")

	return &http.Response{
		StatusCode:    r.status,
		ContentLength: int64(len(data)),
		Body:          ioutil.NopCloser(bytes.NewBuffer(data)),
		Header:        header,
		Request:       req,
	}, nil
}

func Test_RoundTripFunc(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	f := RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		s, err := ioutil.ReadAll(req.Body)
		a.NoError(err)
		a.Equal("this is the request", string(s))
		a.Equal(http.MethodGet, req.Method)
		a.Equal("/foo", req.URL.String())
		header := make(http.Header)
		header.Add("content-type", "text/html; charset=utf-8")
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString("this is the response")),
			Header:     header,
			Request:    req,
		}, nil
	})
	req, err := http.NewRequest(http.MethodGet, "/foo", bytes.NewBufferString("this is the request"))
	a.NoError(err)
	a.NotNil(req)
	res, err := f.RoundTrip(req)
	a.NoError(err)
	a.NotNil(res)
	s, err := ioutil.ReadAll(res.Body)
	a.NoError(err)
	a.Equal("this is the response", string(s))
}
