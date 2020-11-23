package wta_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/wta"
	"github.com/bzimmer/httpwares"
)

func newClient(status int, filename string) (*wta.Client, error) {
	return wta.NewClient(
		wta.WithTransport(&httpwares.TestDataTransport{
			Status:      status,
			Filename:    filename,
			ContentType: "text/html; charset=utf-8",
		}),
	)
}

func newTestRouter(c *wta.Client) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/regions/", wta.RegionsHandler(c))
	r.GET("/reports/", wta.TripReportsHandler(c))
	r.GET("/reports/:reporter", wta.TripReportsHandler(c))
	return r
}

func Test_GetTripReports(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	client, err := newClient(http.StatusOK, "wta_test.html")
	a.NoError(err)
	reports, err := client.Reports.TripReports(context.Background(), "foobar")
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

	c, _ := newClient(http.StatusOK, "wta_test.html")
	r := newTestRouter(c)
	a.NotNil(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/reports/bzimmer", nil)
	r.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)

	var reports wta.TripReports
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&reports)
	a.NoError(err)
	a.NotNil(reports)
	a.NotNil(reports.Reports)
	a.Equal(14, len(reports.Reports))

	// test a response with no html
	c, _ = newClient(http.StatusOK, "wta_test.json")
	r = newTestRouter(c)
	a.NotNil(r)
	w = httptest.NewRecorder()
	req, _ = http.NewRequestWithContext(context.TODO(), http.MethodGet, "/reports/foobar", nil)
	r.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)

	decoder = json.NewDecoder(w.Body)
	err = decoder.Decode(&reports)
	a.NoError(err)
	a.NotNil(reports)
	a.NotNil(reports.Reports)
	a.Equal(0, len(reports.Reports))

	// test with 404 from source
	c, _ = newClient(http.StatusNotFound, "")
	r = newTestRouter(c)
	a.NotNil(r)

	w = httptest.NewRecorder()
	req, _ = http.NewRequestWithContext(context.TODO(), http.MethodGet, "/reports/bzimmer", nil)
	r.ServeHTTP(w, req)
	a.Equal(http.StatusInternalServerError, w.Code)
}

func Test_RegionsHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, _ := newClient(http.StatusOK, "")
	r := newTestRouter(c)
	a.NotNil(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/regions/", nil)
	r.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)

	var regions []wta.Region
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
