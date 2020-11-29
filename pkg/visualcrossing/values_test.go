package visualcrossing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MakeValues(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	v := ForecastOptions{
		Astronomy:      true,
		Units:          UnitsUK,
		AggregateHours: 12,
		Coordinates: Coordinates{
			Latitude:  48.9201,
			Longitude: -122.092,
		},
		AlertLevel: AlertLevelDetail,
	}
	q, err := v.values()
	a.Equal("aggregateHours=12&alertLevel=detail&includeAstronomy=true&locations=48.9201%2C-122.0920&unitGroup=uk", q.Encode())
	a.NoError(err)

	v.AggregateHours = 18
	q, err = v.values()
	a.Nil(q)
	a.Error(err)
}
