package gravl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestFlatten(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var cmds = []*cli.Command{
		{Name: "1", Subcommands: []*cli.Command{{Name: "1a"}}},
		{Name: "2", Subcommands: []*cli.Command{{Name: "2a", Subcommands: []*cli.Command{{Name: "2b"}}}}},
	}
	a.Equal(2, len(cmds))
	a.Equal(5, len(flatten(cmds)))
}
