package store

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/analysis/store/bunt"
	"github.com/bzimmer/gravl/pkg/analysis/store/dynamo"
	"github.com/bzimmer/gravl/pkg/analysis/store/file"
	"github.com/bzimmer/gravl/pkg/analysis/store/strava"
	stravacmd "github.com/bzimmer/gravl/pkg/commands/activity/strava"
)

const DefaultLocalStore = "bunt"

type opener func(*cli.Context, *url.URL) (store.Store, error)

var openers = map[string]opener{
	"file":   openfile,
	"bunt":   openbunt,
	"dynamo": opendynamo,
	"strava": openstrava,
}

func localfile(u *url.URL) string {
	if u.Host != "" {
		return filepath.Join(u.Host, u.Path)
	}
	return u.Path
}

func parse(q string) (*url.URL, error) {
	u, err := url.Parse(q)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" {
		// url.Parse requires `://` so add if missing
		q = fmt.Sprintf("%s://", q)
		u, err = url.Parse(q)
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}

func Open(c *cli.Context, flag string) (store.Store, error) {
	q := c.String(flag)
	if q == "" {
		return nil, errors.New("store not specified")
	}
	u, err := parse(q)
	if err != nil {
		return nil, err
	}
	opener, ok := openers[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("unknown scheme {%s}", u.Scheme)
	}
	return opener(c, u)
}

func opendynamo(c *cli.Context, u *url.URL) (store.Store, error) {
	// @todo(bzimmer) use cli flags for credentials
	// https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/
	values, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	var opts []func(o *config.LoadOptions) error
	for key, vals := range values {
		switch key {
		case "region":
			opts = append(opts, config.WithRegion(vals[0]))
		default:
			return nil, fmt.Errorf("unknown configuration value {%s}", key)
		}
	}
	cfg, err := config.LoadDefaultConfig(c.Context, opts...)
	if err != nil {
		return nil, err
	}
	// @todo(bzimmer) more sensible defaults
	if cfg.Region == "" {
		cfg.Region = "us-west-2"
	}
	return dynamo.Open(c.Context, cfg)
}

func openfile(c *cli.Context, u *url.URL) (store.Store, error) {
	db := localfile(u)
	if db == "" {
		db = c.Args().First()
	}
	return file.Open(db, file.Flush(false))
}

func openstrava(c *cli.Context, u *url.URL) (store.Store, error) {
	client, err := stravacmd.NewAPIClient(c)
	if err != nil {
		return nil, err
	}
	return strava.Open(client), nil
}

func openbunt(c *cli.Context, u *url.URL) (store.Store, error) {
	db := localfile(u)
	if db == "" {
		db = path.Join(xdg.DataHome, pkg.PackageName, "gravl.db")
	}
	return bunt.Open(db)
}
