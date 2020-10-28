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
	a.True(true)

	reader, err := os.Open("testdata/98110_forecast_array.json")
	a.NoError(err)
	decoder := json.NewDecoder(reader)

	fcst := &vc.Forecast{}
	err = decoder.Decode(fcst)
	a.NoError(err)
	a.Equal(1, fcst.QueryCost)
}
