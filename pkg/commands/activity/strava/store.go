package strava

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/analysis/eval/antonmedv"
	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/analysis/store/source/bolt"
	"github.com/bzimmer/gravl/pkg/analysis/store/source/file"
	api "github.com/bzimmer/gravl/pkg/analysis/store/source/strava"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
)

func source(c *cli.Context) (store.Source, error) {
	switch {
	case c.NArg() == 1:
		return file.Open(c.Args().First())
	default:
		client, err := NewAPIClient(c)
		if err != nil {
			return nil, err
		}
		return api.Open(client), nil
	}
}

func sink(c *cli.Context) (store.SourceSink, error) {
	path := c.Path("store")
	if path == "" {
		return nil, errors.New("nil db path")
	}
	return bolt.Open(path)
}

func remove(c *cli.Context) error {
	db, err := sink(c)
	if err != nil {
		return err
	}
	defer db.Close()

	acts, err := db.Activities(c.Context)
	if err != nil {
		return err
	}
	if c.IsSet("filter") {
		evaluator := antonmedv.New()
		acts, err = evaluator.Filter(c.Context, c.String("filter"), acts)
		if err != nil {
			return err
		}
	}
	ids := make([]int64, len(acts))
	for i, act := range acts {
		ids[i] = act.ID
	}
	if c.Bool("dryrun") {
		return encoding.Encode(ids)
	}
	if err := db.Remove(c.Context, acts...); err != nil {
		return err
	}
	return encoding.Encode(ids)
}

func update(c *cli.Context) error {
	db, err := sink(c)
	if err != nil {
		return err
	}
	defer db.Close()
	src, err := source(c)
	if err != nil {
		return err
	}
	defer src.Close()
	acts, err := src.Activities(c.Context)
	if err != nil {
		return err
	}
	var n int
	for i := 0; i < len(acts); i++ {
		var ok bool
		act := acts[i]
		ok, err = db.Exists(c.Context, act.ID)
		if err != nil {
			return err
		}
		if ok {
			continue
		}
		log.Info().Int64("ID", act.ID).Msg("querying activity details")
		act, err = src.Activity(c.Context, act.ID)
		if err != nil {
			return err
		}
		log.Info().Int("n", n+1).Int64("ID", act.ID).Str("name", act.Name).Msg("saving activity details")
		if err = db.Save(c.Context, act); err != nil {
			return nil
		}
		n++
	}
	if err = encoding.Encode(map[string]int{
		"total":    len(acts),
		"new":      n,
		"existing": len(acts) - n}); err != nil {
		return err
	}
	return nil
}

var updateCommand = &cli.Command{
	Name:   "update",
	Usage:  "Query and update Strava activities to local storage",
	Action: update,
}

var removeCommand = &cli.Command{
	Name:  "remove",
	Usage: "Remove activities from local storage",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "filter",
			Aliases:  []string{"f"},
			Required: true,
			Usage:    "Expression for filtering activities to remove",
		},
		&cli.BoolFlag{
			Name:    "dryrun",
			Aliases: []string{"n"},
			Value:   false,
			Usage:   "Don't actually remove anything, just show what would be done",
		},
	},
	Action: remove,
}

var storeCommand = &cli.Command{
	Name:  "store",
	Usage: "Manage a local store of Strava activities",
	Flags: []cli.Flag{commands.StoreFlag},
	Subcommands: []*cli.Command{
		updateCommand,
		removeCommand,
	},
}
