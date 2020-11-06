package common

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen)
	yellow = color.New(color.FgYellow)
)

// VerboseTransport .
type VerboseTransport struct {
	Transport http.RoundTripper
}

func (t *VerboseTransport) isText(header http.Header) bool {
	contentType := header.Get("Content-Type")
	if contentType == "" {
		return false
	}
	if strings.HasPrefix(contentType, "text/") {
		return true
	}
	// content-type is two parts:
	//  - type
	//  - parameters
	splits := strings.Split(contentType, ";")
	switch splits[0] {
	case "application/json":
	case "application/ld+json":
	case "application/geojson":
	default:
		return false
	}
	return true
}

// RoundTrip .
func (t *VerboseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := http.DefaultTransport
	if t.Transport != nil {
		transport = t.Transport
	}
	dump, _ := httputil.DumpRequestOut(req, t.isText(req.Header))
	yellow.Fprintln(color.Error, string(dump))
	res, err := transport.RoundTrip(req)
	dump, _ = httputil.DumpResponse(res, t.isText(res.Header))
	green.Fprintln(color.Error, string(dump))
	return res, err
}

// Requester .
type Requester func(*http.Request) error

// Responder .
type Responder func(*http.Response) error

// TestDataTransport .
type TestDataTransport struct {
	Status      int
	Filename    string
	ContentType string
	Requester   Requester
	Responder   Responder
}

// RoundTrip .
func (t *TestDataTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		err  error
		data []byte
	)
	if t.Requester != nil {
		err = t.Requester(req)
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
	header.Set("Content-Type", t.ContentType)

	res := &http.Response{
		StatusCode:    t.Status,
		ContentLength: int64(len(data)),
		Body:          ioutil.NopCloser(bytes.NewBuffer(data)),
		Header:        header,
		Request:       req,
	}
	if t.Responder != nil {
		err = t.Responder(res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
