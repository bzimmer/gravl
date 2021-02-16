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

	var x int
	var y string
	var z bool

	opt, err := options.Parse("database,X=3")
	a.NoError(err)

	a.Equal("database", opt.Name)
	a.Equal(1, len(opt.Options))
	a.Equal("3", opt.Options["X"])

	fs := flag.NewFlagSet("test", flag.ExitOnError)
	fs.IntVar(&x, "X", x, "number of Xs")
	fs.StringVar(&y, "Y", y, "Ys")
	fs.BoolVar(&z, "Z", z, "use z?")

	a.NoError(opt.ApplyFlags(fs))
	a.Equal(3, x)

	opt, err = options.Parse("database,X=10,Z")
	a.Error(err)
	a.Nil(opt)

	x, y, z = 0, "", false
	opt, err = options.Parse("database,X=10,Y=hello")
	a.NoError(err)
	a.NotNil(opt)
	a.NoError(opt.ApplyFlags(fs))
	a.Equal(10, x)
	a.Equal("hello", y)

	x, y, z = 0, "", false
	opt, err = options.Parse("database,Z=true,X=10")
	a.NoError(err)
	a.NotNil(opt)
	a.NoError(opt.ApplyFlags(fs))
	a.Equal(10, x)
	a.True(z)
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
