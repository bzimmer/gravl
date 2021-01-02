package benford_test

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/analysis/passes/benford"
)

func TestBenfordsLaw(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	file, err := os.Open("testdata/alps.csv")
	a.NoError(err)
	defer file.Close()

	r := csv.NewReader(file)
	r.LazyQuotes = true
	_, _ = r.Read()

	var elevations []int
	for {
		x, _ := r.Read()
		if x == nil {
			break
		}
		elv, err := strconv.Atoi(x[2])
		a.NoError(err)
		elevations = append(elevations, elv)
	}
	b := benford.Law(elevations)
	a.InEpsilon(0.79663496, b.ChiSquared, 0.0001)
}
