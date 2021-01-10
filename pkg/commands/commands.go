package commands

import (
	"crypto/rand"
	"encoding/base64"
	"path"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

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

func Token(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

var StoreFlag = func() cli.Flag {
	store := path.Join(xdg.DataHome, pkg.PackageName, "gravl.db")
	return altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "store",
		Value: store,
		Usage: "Path to the database",
	})
}()
