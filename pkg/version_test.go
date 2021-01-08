package pkg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	a.Equal("development", pkg.BuildVersion)
}
