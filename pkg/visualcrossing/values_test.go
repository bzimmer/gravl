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
		WithLocation("48.9201,-122.092"),
		WithAlerts(AlertLevelDetail),
	})
	a.NoError(err)
	a.NotNil(v)
	q := v.Encode()
	a.Equal("aggregateHours=12&alertLevel=detail&includeAstronomy=true&locations=48.9201%2C-122.092&unitGroup=uk", q)

	// test re-using options
	v, err = makeValues([]ForecastOption{
		WithUnits(UnitsUS),
		WithUnits(UnitsUK),
		WithUnits(UnitsMetric),
	})
	q = v.Encode()
	a.NoError(err)
	a.Equal("unitGroup=metric", q)

	v = &url.Values{}
	a.Error(WithUnits("foo")(v))
	a.Error(WithAlerts("bar")(v))
	a.Error(WithAggregateHours(-1000)(v))
}