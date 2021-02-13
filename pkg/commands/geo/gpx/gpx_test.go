package gpx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/bzimmer/gravl/pkg/commands/internal"
)

func TestInfoIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	a := assert.New(t)

	c := internal.Gravl("gpx", "info", "testdata/2017-07-13-TdF-Stage18.gpx")
	<-c.Start()
	a.True(c.Success())

	val := gjson.Parse(c.Stdout())
	a.InEpsilon(2310, val.Get("ascent").Num, 1.0)
}
