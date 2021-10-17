//go:build integration

package qp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/internal"
)

func TestQPIntegration(t *testing.T) {
	a := assert.New(t)
	c := internal.Gravl("-c", "qp", "-e", "zwift", "-u", "ca")
	<-c.Start()
	a.True(c.Success())
}
