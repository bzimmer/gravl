package store

import (
	"errors"
	"fmt"
	"path"
	"strconv"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
	stravacmd "github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/options"
	"github.com/bzimmer/gravl/pkg/store"
	"github.com/bzimmer/gravl/pkg/store/bunt"
	"github.com/bzimmer/gravl/pkg/store/fake"
	"github.com/bzimmer/gravl/pkg/store/file"
	"github.com/bzimmer/gravl/pkg/store/strava"
)

const DefaultLocalStore = "bunt"

type opener func(*cli.Context, *options.Option) (store.Store, error)

var openers = map[string]opener{
	"file":   openfile,
	"bunt":   openbunt,
	"strava": openstrava,
	"fake":   openfake,
}

func Open(c *cli.Context, flag string) (store.Store, error) {
	q := c.String(flag)
	if q == "" {
		return nil, errors.New("store not specified")
	}
	u, err := options.Parse(q)
	if err != nil {
		return nil, err
	}
	opener, ok := openers[u.Name]
	if !ok {
		return nil, fmt.Errorf("unknown scheme {%s}", u.Name)
	}
	return opener(c, u)
}

func openfile(c *cli.Context, u *options.Option) (store.Store, error) {
	db, ok := u.Options["path"]
	if !ok {
		return nil, errors.New("missing filename")
	}
	return file.Open(db, file.Flush(false))
}

func openstrava(c *cli.Context, u *options.Option) (store.Store, error) {
	client, err := stravacmd.NewAPIClient(c)
	if err != nil {
		return nil, err
	}
	return strava.Open(client), nil
}

func openbunt(c *cli.Context, u *options.Option) (store.Store, error) {
	db, ok := u.Options["path"]
	if !ok {
		db = path.Join(xdg.DataHome, pkg.PackageName, "gravl.db")
	}
	return bunt.Open(db)
}

func openfake(c *cli.Context, u *options.Option) (store.Store, error) {
	n := 100
	x, ok := u.Options["n"]
	if ok {
		y, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return nil, err
		}
		n = int(y)
	}
	fuzz := false
	x, ok = u.Options["fuzz"]
	if ok {
		y, err := strconv.ParseBool(x)
		if err != nil {
			return nil, err
		}
		fuzz = y
	}
	return fake.Open(n, fuzz)
}
