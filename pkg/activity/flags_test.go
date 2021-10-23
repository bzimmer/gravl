package activity_test

import (
	"testing"

	"github.com/bzimmer/gravl/pkg/activity"
	"github.com/stretchr/testify/assert"
)

func TestFlags(t *testing.T) {
	a := assert.New(t)
	flags := activity.RateLimitFlags()
	a.Equal(3, len(flags))
}
