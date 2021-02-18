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
		x      int
		y      string
		z      bool
		num    int
		name   string
		option string
		err    bool
	}{
		{option: "", err: true},
		{option: "database,X=10,Z", err: true},
		{option: "X=10", err: false, name: "X=10"},
		{option: "database,X=3", name: "database", x: 3, num: 1},
		{option: "database,X=10,Y=hello", name: "database", x: 10, y: "hello", num: 2},
		{option: "database,Z=true,X=10", name: "database", x: 10, num: 2, z: true},
		{option: "database,Z=true,X=10", name: "database", x: 10, num: 2, z: true},
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
			a.Equal(v.name, opt.Name)
			a.Equal(v.num, len(opt.Options))

			a.NoError(opt.ApplyFlags(fs))
			a.Equal(v.x, x)
			a.Equal(v.y, y)
			a.Equal(v.z, z)
		})
	}
}

func TestOptionsNoEquals(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	opt, err := options.Parse("database,file:/home/somebody/gravl.db")
	a.Error(err)
	a.Nil(opt)

	opt, err = options.Parse("database")
	a.NoError(err)
	a.Equal("database", opt.Name)
	a.Equal(0, len(opt.Options))

	opt, err = options.Parse("")
	a.Error(err)
	a.Nil(opt)
}
