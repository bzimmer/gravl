package version_test

import (
	"testing"

	"github.com/bzimmer/gravl/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	a.Equal("development", version.BuildVersion)
}
