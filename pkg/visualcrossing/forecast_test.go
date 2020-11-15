package visualcrossing

import (
	"encoding/json"
	"net/url"
	"os"
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

func Test_Model(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	reader, err := os.Open("testdata/forecast.json")
	a.NoError(err)
	a.NotNil(reader)

	var fcst forecast
	err = json.NewDecoder(reader).Decode(&fcst)
	a.NoError(err)

	a.Equal(1, fcst.QueryCost)
	a.Equal(1, len(fcst.Locations))

	loc := fcst.Locations[0]
	fc := loc.ForecastConditions
	a.Equal(16, len(fc))

	cond := fc[len(fc)-1]
	a.Equal(32.1, cond.WindChill)
}
