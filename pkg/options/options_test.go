package options_test

import (
	"flag"
	"testing"

	"github.com/bzimmer/gravl/pkg/options"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	tests := []struct {
		x, num          int
		y, name, option string
		z, err          bool
	}{
		{option: "", err: true},
		{option: "database", name: "database"},
		{option: "database,X=10,Z", err: true},
		{option: "X=10", name: "X=10"},
		{option: "database,X=3", name: "database", x: 3, num: 1},
		{option: "database,X=10,Y=hello", name: "database", x: 10, y: "hello", num: 2},
		{option: "database,file:/home/somebody/gravl.db", err: true},
		{option: "database,Z=true,X=10", name: "database", x: 10, num: 2, z: true},
		{option: "database,Z=true,X=10", name: "database", x: 10, num: 2, z: true},

		// @todo(bzimmer)
		// {option: `database,X=10,Y="hello, goodbye"`, name: "database", x: 10, y: "hello, goodbye", num: 2},
	}
	for _, tt := range tests {
		v := tt
		t.Run(v.option, func(t *testing.T) {
			var x int
			var y string
			var z bool
			fs := flag.NewFlagSet("test", flag.ExitOnError)
			fs.IntVar(&x, "X", x, "number of Xs")
			fs.StringVar(&y, "Y", y, "Ys")
			fs.BoolVar(&z, "Z", z, "use z?")

			opt, err := options.Parse(v.option)
			if v.err {
				a.Error(err)
				return
			}
			a.NoError(err)
			a.NotNil(opt)
			if opt == nil {
				return
			}
			a.Equal(v.name, opt.Name)
			a.Equal(v.num, len(opt.Options))

			a.NoError(opt.ApplyFlags(fs))
			a.Equal(v.x, x)
			a.Equal(v.y, y)
			a.Equal(v.z, z)
		})
	}
}
