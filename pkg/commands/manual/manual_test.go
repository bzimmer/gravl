package manual_test

import (
	"testing"

	"github.com/bzimmer/gravl/pkg/commands/internal"
	"github.com/stretchr/testify/assert"
)

func TestManualIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	a := assert.New(t)

	tests := []struct {
		err  bool
		name string
		args []string
	}{
		{name: "manual no flags", err: true, args: []string{"manual"}},
		{name: "manual both args", err: true, args: []string{"manual", "--commands", "--analyzers"}},
		{name: "manual commands", args: []string{"manual", "--commands"}},
		{name: "manual analyzers", args: []string{"manual", "--analyzers"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := internal.Gravl(tt.args...)
			<-c.Start()
			switch tt.err {
			case true:
				a.False(c.Success())
				a.Equal(0, len(c.Stdout()))
			case false:
				a.True(c.Success())
				a.Greater(len(c.Stdout()), 0)
			}
		})
	}
}
