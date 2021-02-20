package wta_test

import (
	"testing"

	"github.com/bzimmer/gravl/pkg/commands/internal"
	"github.com/stretchr/testify/assert"
)

func TestReportsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test suite")
		return
	}
	t.Parallel()
	a := assert.New(t)

	tests := []struct {
		name, reporter string
	}{
		{name: "recent", reporter: ""},
		{name: "bzimmer", reporter: "bzimmer"},
	}
	for _, tt := range tests {
		v := tt
		t.Run(v.name, func(t *testing.T) {
			c := internal.Gravl("wta", v.name)
			<-c.Start()
			a.True(c.Success())
		})
	}
}
