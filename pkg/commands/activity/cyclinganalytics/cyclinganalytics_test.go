package cyclinganalytics_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/bzimmer/gravl/pkg/commands/internal"
)

const N = 22

// @todo(bzimmer) enable geojson support for cyclinganalytics
const geojson = false

func TestAthleteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	a := assert.New(t)

	c := internal.Gravl("-c", "cyclinganalytics", "athlete")
	<-c.Start()
	a.True(c.Success())
}

func TestActivityIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	a := assert.New(t)

	c := internal.Gravl("-c", "cyclinganalytics", "activities", "-N", strconv.FormatInt(N, 10))
	<-c.Start()
	a.True(c.Success())

	var i int
	var randomID = rand.Intn(N - 1) // nolint
	gjson.ForEachLine(c.Stdout(), func(res gjson.Result) bool {
		id := gjson.Get(res.String(), "id").Int()
		a.Greater(id, int64(0))
		if i == randomID {
			idS := strconv.FormatInt(id, 10)
			c = internal.Gravl("-c", "cyclinganalytics", "activity", idS)
			<-c.Start()
			a.True(c.Success())
			c = internal.Gravl("-e", "gpx", "cyclinganalytics", "activity", idS)
			<-c.Start()
			a.True(c.Success())
			if geojson {
				c = internal.Gravl("-e", "geojson", "cyclinganalytics", "activity", idS)
				<-c.Start()
				a.True(c.Success())
				res = gjson.Parse(c.Stdout())
				a.NotNil(res)
			}
			c = internal.Gravl("-c", "cyclinganalytics", "activity", idS)
			<-c.Start()
			a.True(c.Success())
		}
		i++
		return true
	})
	a.Equal(N, i)
}