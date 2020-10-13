package wta

import (
	"bytes"
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

type TestDataTransport struct{}

func (r *TestDataTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	dir, _ := os.Getwd()
	data, _ := ioutil.ReadFile(filepath.Join(dir, "wta_test.html"))

	header := make(http.Header)
	header.Add("content-type", "text/html; charset=utf-8")

	return &http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		ContentLength: int64(len(data)),
		Body:          ioutil.NopCloser(bytes.NewBuffer(data)),
		Header:        header,
		Request:       req,
	}, nil
}

func newCollector() *colly.Collector {
	c := colly.NewCollector()
	c.WithTransport(&TestDataTransport{})
	return c
}

func newTestRouter() (*gin.Engine, error) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/regions/", RegionsHandler())
	r.GET("/reports/:reporter", TripReportsHandler(newCollector()))
	return r, nil
}

func Test_Query(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	q, err := Query("foobar")
	a.NoError(err)
	a.NotNil(q)
}

func Test_GetTripReports(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	q, err := Query("foobar")
	a.NoError(err)

	reports, err := GetTripReports(newCollector(), q.String())
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

	r, err := newTestRouter()
	a.NoError(err)
	a.NotNil(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/reports/bzimmer", nil)
	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	var reports TripReports
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&reports)
	a.NoError(err)
	a.NotNil(reports)
	a.NotNil(reports.Reports)
	a.Equal(14, len(reports.Reports))
}

func Test_RegionsHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	r, err := newTestRouter()
	a.NoError(err)
	a.NotNil(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/regions/", nil)
	r.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	var regions []Region
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&regions)
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
