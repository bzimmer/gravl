package wx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/providers/wx"
)

func TestWindBearing(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	b, err := wx.WindBearing("W")
	a.NoError(err)
	a.Equal(270.0, b)
	b, err = wx.WindBearing("bar")
	a.Zero(b)
	a.Error(err)
	b, err = wx.WindBearing("")
	a.Zero(b)
	a.Error(err)
}

func TestCompassPoint(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	for _, x := range []float64{358.0, 3.87, 360.0, 354.375, 0.000, 5.624} {
		p, err := wx.CompassPoint(x)
		a.NoError(err)
		a.Equal("N", p)
	}
	p, err := wx.CompassPoint(238.3)
	a.NoError(err)
	a.Equal("SWbW", p)
	p, err = wx.CompassPoint(363.2)
	a.Error(err)
	a.Equal("", p)
}

func TestRoundTrip(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	c, err := wx.CompassPoint(160.2)
	a.NoError(err)
	a.Equal("SSE", c)
	b, err := wx.WindBearing(c)
	a.NoError(err)
	a.Equal(157.50, b)
}
