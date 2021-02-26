package qp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/commands/internal"
)

func TestQPIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test suite")
		return
	}
	t.Parallel()
	a := assert.New(t)
	c := internal.Gravl("-c", "qp", "-e", "zwift", "-u", "ca")
	<-c.Start()
	a.True(c.Success())
}
