package commands

import (
	"path"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
)

func Merge(flags ...[]cli.Flag) []cli.Flag {
	var f []cli.Flag
	for _, x := range flags {
		f = append(f, x...)
	}
	return f
}

func Before(befores ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		for _, fn := range befores {
			if fn == nil {
				continue
			}
			if e := fn(c); e != nil {
				return e
			}
		}
		return nil
	}
}

var StoreFlag = func() cli.Flag {
	store := path.Join(xdg.DataHome, pkg.PackageName, "gravl.db")
	return &cli.PathFlag{
		Name:      "store",
		Value:     store,
		TakesFile: true,
		Usage:     "Path to the database",
	}
}()
