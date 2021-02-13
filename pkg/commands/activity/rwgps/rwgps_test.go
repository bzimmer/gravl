package rwgps_test

import (
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/valyala/fastjson"

	"github.com/bzimmer/gravl/pkg/commands/internal"
)

func TestRWGPSAthleteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	a := assert.New(t)

	c := testcli.Command(internal.PackageGravl(), "-c", "rwgps", "athlete")
	c.Run()
	a.True(c.Success())
}
func TestRWGPSActivityIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	a := assert.New(t)

	c := testcli.Command(internal.PackageGravl(), "-c", "rwgps", "activities", "-N", "25")
	c.Run()
	a.True(c.Success())

	var sc fastjson.Scanner
	sc.InitBytes([]byte(c.Stdout()))

	var i int
	var line string
	for ; sc.Next(); i++ {
		line = sc.Value().String()
	}

	a.Equal(25, i)
	res := gjson.Get(line, "id")
	a.Greater(res.Int(), int64(0))

	c = testcli.Command(internal.PackageGravl(), "-c", "rwgps", "activity", res.String())
	c.Run()
	a.True(c.Success())
}
