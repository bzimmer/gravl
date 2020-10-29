package visualcrossing_test

import (
	"encoding/json"
	"os"
	"testing"

	vc "github.com/bzimmer/wta/pkg/visualcrossing"

	"github.com/stretchr/testify/assert"
)

func Test_Model(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	reader, err := os.Open("testdata/98110_forecast_array.json")
	a.NoError(err)
	decoder := json.NewDecoder(reader)

	fcst := &vc.Forecast{}
	err = decoder.Decode(fcst)
	a.NoError(err)

	a.Equal(1, fcst.QueryCost)
	a.Equal(1, len(fcst.Locations))

	loc := fcst.Locations[0]
	a.Equal(16, len(loc.Conditions))

	cond := loc.Conditions[len(loc.Conditions)-1]
	a.Equal(32.1, cond.WindChill)
}
