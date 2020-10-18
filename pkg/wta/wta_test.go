package wta

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
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

func newCollectorWithFilename(filename string) *colly.Collector {
	return newCollector(http.StatusOK, filename)
}

func newCollectorWithStatus(status int) *colly.Collector {
	return newCollector(status, "")
}

func newCollector(status int, filename string) *colly.Collector {
	c := NewCollector()
	c.WithTransport(&TestDataTransport{
		status:   status,
		filename: filename,
	})
	return c
}

func newClient(c *colly.Collector) (*Client, error) {
	return NewClient(WithCollector(c))
}

func newTestRouter(c *colly.Collector) (*gin.Engine, error) {
	gin.SetMode(gin.TestMode)
	client, err := newClient(c)
	if err != nil {
		return nil, err
	}
	r := NewRouter(client)
	return r, nil
}

func Test_query(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	q := query("foobar")
	a.NotNil(q)
}

func Test_GetTripReports(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	q := query("foobar")
	c := newCollectorWithFilename("wta_test.html")
	client, err := newClient(c)
	reports, err := client.Reports.TripReports(context.Background(), q.String())
	a.NoError(err)
	a.Equal(14, len(reports))

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].HikeDate.After(reports[j].HikeDate)
	})

	tr := reports[5]
	a.Equal(4, tr.Photos)
	a.Equal(13, tr.Votes)
	a.Equal("Lake Angeles, Klahhane Ridge", tr.Title)
}

func Test_TripReportsHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	// test with known good html
	c := newCollectorWithFilename("wta_test.html")
	r, _ := newTestRouter(c)
	a.NotNil(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/reports/bzimmer", nil)
	r.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)

	var reports TripReports
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&reports)
	a.NoError(err)
	a.NotNil(reports)
	a.NotNil(reports.Reports)
	a.Equal(14, len(reports.Reports))

	// test a response with no html
	c = newCollectorWithFilename("wta_test.json")
	r, _ = newTestRouter(c)
	a.NotNil(r)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/reports/foobar", nil)
	r.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)

	decoder = json.NewDecoder(w.Body)
	err = decoder.Decode(&reports)
	a.NoError(err)
	a.NotNil(reports)
	a.NotNil(reports.Reports)
	a.Equal(0, len(reports.Reports))

	// test with 404 from source
	c = newCollectorWithStatus(http.StatusNotFound)
	r, _ = newTestRouter(c)
	a.NotNil(r)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/reports/bzimmer", nil)
	r.ServeHTTP(w, req)
	a.Equal(http.StatusInternalServerError, w.Code)
}

func Test_RegionsHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c := newCollectorWithFilename("wta_test.html")
	r, _ := newTestRouter(c)
	a.NotNil(r)
	a.NotNil(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/regions/", nil)
	r.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)

	var regions []Region
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&regions)
	a.NoError(err)
	a.NotNil(regions)
	a.NotNil(regions)
	a.Equal(11, len(regions))

	var id string
	for _, region := range regions {
		if region.Name == "Olympic Peninsula" {
			id = region.ID
		}
	}
	a.Equal("922e688d784aa95dfb80047d2d79dcf6", id)
}

func Test_VersionHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	// data never used
	c := newCollectorWithFilename("wta_test.html")
	r, _ := newTestRouter(c)
	a.NotNil(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/version/", nil)
	r.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)

	var version map[string]string
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&version)
	a.NoError(err)
	a.Equal("development", version["build_version"])
}
