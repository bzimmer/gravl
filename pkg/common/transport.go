package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"path/filepath"
)

// VerboseTransport .
type VerboseTransport struct {
	Transport http.RoundTripper
}

// RoundTrip .
func (t *VerboseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := http.DefaultTransport
	if t.Transport != nil {
		transport = t.Transport
	}
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(string(dump))
	res, err := transport.RoundTrip(req)
	dump, _ = httputil.DumpResponse(res, true)
	fmt.Println(string(dump))
	return res, err
}

// TestDataTransport .
type TestDataTransport struct {
	Status      int
	Filename    string
	ContentType string
	Validator   func(*http.Request) error
}

// RoundTrip .
func (t *TestDataTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		err  error
		data []byte
	)
	if t.Validator != nil {
		err = t.Validator(req)
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
