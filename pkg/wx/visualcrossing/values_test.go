package visualcrossing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twpayne/go-geom"
)

func Test_MakeValues(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	v := ForecastOptions{
		Astronomy:      true,
		Units:          UnitsUK,
		AggregateHours: 12,
		Point:          geom.NewPointFlat(geom.XY, []float64{-122.092, 48.9201}),
		AlertLevel:     AlertLevelDetail,
	}
	q, err := v.values()
	a.Equal("aggregateHours=12&alertLevel=detail&includeAstronomy=true&locations=48.9201%2C-122.0920&unitGroup=uk", q.Encode())
	a.NoError(err)

	v.AggregateHours = 18
	q, err = v.values()
	a.Nil(q)
	a.Error(err)

	v = ForecastOptions{
		Astronomy:      true,
		Units:          UnitsUK,
		AggregateHours: 12,
		Location:       "Seattle, WA",
		Point:          geom.NewPointFlat(geom.XY, []float64{-122.092, 48.9201}),
		AlertLevel:     AlertLevelDetail,
	}
	q, err = v.values()
	a.Equal("aggregateHours=12&alertLevel=detail&includeAstronomy=true&locations=Seattle%2C+WA&unitGroup=uk", q.Encode())
	a.NoError(err)

	v = ForecastOptions{Location: "Seattle, WA"}
	q, err = v.values()
	a.Equal("alertLevel=none&includeAstronomy=false&locations=Seattle%2C+WA&unitGroup=metric", q.Encode())
	a.NoError(err)
}
