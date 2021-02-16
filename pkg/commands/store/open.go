package store

import (
	"errors"
	"fmt"
	"path"

	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
	stravacmd "github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/options"
	"github.com/bzimmer/gravl/pkg/store"
	"github.com/bzimmer/gravl/pkg/store/bunt"
	"github.com/bzimmer/gravl/pkg/store/file"
	"github.com/bzimmer/gravl/pkg/store/strava"
)

const DefaultLocalStore = "bunt"

type opener func(*cli.Context, *options.Option) (store.Store, error)

var openers = map[string]opener{
	"file":   openfile,
	"bunt":   openbunt,
	"strava": openstrava,
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
	log.Info().Str("path", db).Msg("file db")
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
	log.Info().Str("path", db).Msg("bunt db")
	return bunt.Open(db)
}

// func opendynamo(c *cli.Context, u *options.Option) (store.Store, error) {
// 	// @todo(bzimmer) use cli flags for credentials
// 	// https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/
// 	var opts []func(o *config.LoadOptions) error
// 	for key, vals := range u.Options {
// 		switch key {
// 		case "region":
// 			opts = append(opts, config.WithRegion(vals))
// 		default:
// 			return nil, fmt.Errorf("unknown configuration value {%s}", key)
// 		}
// 	}
// 	cfg, err := config.LoadDefaultConfig(c.Context, opts...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// @todo(bzimmer) more sensible defaults
// 	if cfg.Region == "" {
// 		cfg.Region = "us-west-2"
// 	}
// 	return dynamo.Open(c.Context, cfg)
// }
