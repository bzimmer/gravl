package visualcrossing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twpayne/go-geom"

	"github.com/bzimmer/gravl/pkg/providers/wx"
)

func Test_MakeValues(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	v := wx.ForecastOptions{
		Units:          wx.Metric,
		AggregateHours: 12,
		Point:          geom.NewPointFlat(geom.XY, []float64{-122.092, 48.9201}),
	}
	q, err := values(v)
	a.Equal("aggregateHours=12&alertLevel=detail&includeAstronomy=true&locations=48.9201%2C-122.0920&unitGroup=metric", q.Encode())
	a.NoError(err)

	v.AggregateHours = 18
	q, err = values(v)
	a.Nil(q)
	a.Error(err)

	v = wx.ForecastOptions{
		Units:          wx.Imperial,
		AggregateHours: 12,
		Location:       "Seattle, WA",
		Point:          geom.NewPointFlat(geom.XY, []float64{-122.092, 48.9201}),
	}
	q, err = values(v)
	a.Equal("aggregateHours=12&alertLevel=detail&includeAstronomy=true&locations=Seattle%2C+WA&unitGroup=us", q.Encode())
	a.NoError(err)

	v = wx.ForecastOptions{Location: "Seattle, WA"}
	q, err = values(v)
	a.Equal("alertLevel=detail&includeAstronomy=true&locations=Seattle%2C+WA&unitGroup=metric", q.Encode())
	a.NoError(err)
}
