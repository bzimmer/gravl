package visualcrossing_test

import (
	"testing"

	"github.com/bzimmer/gravl/pkg/visualcrossing"
	"github.com/stretchr/testify/assert"
)

func Test_Enums(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	a.Equal("metric", visualcrossing.UnitsMetric.String())
	a.Equal("summary", visualcrossing.AlertLevelSummary.String())
}
