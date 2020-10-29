package visualcrossing

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MakeValues(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	v, err := makeValues([]ForecastOption{
		WithAstronomy(true),
		WithUnits(UnitsUK),
		WithAggregateHours(12),
		WithLocation("48.9201,-122.092", "Missoula, MT", "Basel, Switzerland"),
		WithAlerts(AlertLevelDetail),
	})
	a.NoError(err)
	a.NotNil(v)

	q := v.Encode()
	a.Equal("aggregateHours=12&alertLevel=detail&includeAstronomy=true&locations=48.9201%2C-122.092%7CMissoula%2C+MT%7CBasel%2C+Switzerland&unitGroup=uk", q)

	v = &url.Values{}
	a.Error(WithUnits("foo")(v))
	a.Error(WithAlerts("bar")(v))
}
