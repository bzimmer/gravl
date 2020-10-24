package rwgps

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestDataTransport struct {
	status      int
	filename    string
	contentType string
}

func (t *TestDataTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var data []byte
	if t.filename != "" {
		dir, _ := os.Getwd()
		filename := filepath.Join(dir, "../../testdata", t.filename)
		data, _ = ioutil.ReadFile(filename)
	} else {
		data = make([]byte, 0)
	}

	header := make(http.Header)
	header.Add("Content-Type", t.contentType)

	return &http.Response{
		StatusCode:    t.status,
		ContentLength: int64(len(data)),
		Body:          ioutil.NopCloser(bytes.NewBuffer(data)),
		Header:        header,
		Request:       req,
	}, nil
}

func newClient(status int, filename string) (*Client, error) {
	return NewClient(
		WithTransport(&TestDataTransport{
			status:      status,
			filename:    filename,
			contentType: "application/json",
		}))
}

func Test_Trip(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "rwgps_trip_94.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	fc, err := c.Trips.Trip(ctx, 94)
	a.NoError(err)
	a.NotNil(fc)
	a.Equal(94, fc.Features[0].ID)
	a.Equal("trip", fc.Features[0].Properties["type"])
	a.True(fc.Features[0].Geometry.IsLineString())
	a.Equal(1465, len(fc.Features[0].Geometry.LineString))

	fc, err = c.Trips.Trip(nil, 94)
	a.Error(err)
	a.Nil(fc)
}

func Test_Route(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := newClient(http.StatusOK, "rwgps_route_141014.json")
	a.NoError(err)
	a.NotNil(c)

	ctx := context.Background()
	fc, err := c.Trips.Route(ctx, 141014)
	a.NoError(err)
	a.NotNil(fc)
	a.Equal(141014, fc.Features[0].ID)
	a.Equal("route", fc.Features[0].Properties["type"])
	a.Equal(1, len(fc.Features))
	a.True(fc.Features[0].Geometry.IsLineString())
	a.Equal(1154, len(fc.Features[0].Geometry.LineString))

	fc, err = c.Trips.Route(nil, 141014)
	a.Error(err)
	a.Nil(fc)
}
