package wta_test

import (
	"context"
	"net/http"
	"sort"
	"testing"

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
