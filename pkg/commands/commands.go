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

// Merge multiple slices of cli flags into a single slice
func Merge(flags ...[]cli.Flag) []cli.Flag {
	var f []cli.Flag
	for _, x := range flags {
		f = append(f, x...)
	}
	return f
}

// Before combines multiple before functions into a single before functions
func Before(befores ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		for _, fn := range befores {
			if fn == nil {
				continue
			}
			if err := fn(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// Token produces a random token of length `n`
func Token(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// StoreFlag for the path local storage
var StoreFlag = func() cli.Flag {
	store := path.Join(xdg.DataHome, pkg.PackageName, "gravl.db")
	return altsrc.NewStringFlag(&cli.StringFlag{
		Name:  "store",
		Value: store,
		Usage: "Path to the database",
	})
}()
