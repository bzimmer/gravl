package store

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/analysis/store/bunt"
	"github.com/bzimmer/gravl/pkg/analysis/store/file"
	stravastore "github.com/bzimmer/gravl/pkg/analysis/store/strava"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	stravaapi "github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func source(c *cli.Context) (store.Source, error) {
	switch {
	case c.NArg() == 1:
		return file.Open(c.Args().First())
	default:
		client, err := strava.NewAPIClient(c)
		if err != nil {
			return nil, err
		}
		return stravastore.Open(client), nil
	}
}

func sink(c *cli.Context) (store.SourceSink, error) {
	path := c.Path("store")
	if path == "" {
		return nil, errors.New("nil db path")
	}
	return bunt.Open(path)
}

func filter(c *cli.Context, acts []*stravaapi.Activity) ([]*stravaapi.Activity, error) {
	if !c.IsSet("filter") {
		return acts, nil
	}
	evaluator := commands.Filterer(c.String("filter"))
	return evaluator.Filter(c.Context, acts)
}

func export(c *cli.Context) error {
	db, err := sink(c)
	if err != nil {
		return err
	}
	defer db.Close()
	ca, ce := db.Activities(c.Context)
	acts, err := stravaapi.Activities(c.Context, ca, ce)
	if err != nil {
		return err
	}
	acts, err = filter(c, acts)
	if err != nil {
		return err
	}
	for _, act := range acts {
		if err := encoding.Encode(act); err != nil {
			return err
		}
	}
	return nil
}

func remove(c *cli.Context) error {
	db, err := sink(c)
	if err != nil {
		return err
	}
	defer db.Close()
	ca, ce := db.Activities(c.Context)
	acts, err := stravaapi.Activities(c.Context, ca, ce)
	if err != nil {
		return err
	}
	acts, err = filter(c, acts)
	if err != nil {
		return err
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
	var ok bool
	var err error
	var total, n int
	output, err := sink(c)
	if err != nil {
		return err
	}
	defer output.Close()
	input, err := source(c)
	if err != nil {
		return err
	}
	defer input.Close()
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()
	acts, errs := input.Activities(ctx)
	for active := true; active; {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err, ok = <-errs:
			// if the channel is not closed, return the error
			// if the channel is closed, do nothing to ensure all activities are consumed
			if ok {
				return err
			}
		case act, ok := <-acts:
			if !ok {
				// break the loop to return the processing results
				active = false
				break
			}
			total++
			ok, err = output.Exists(ctx, act.ID)
			if err != nil {
				return err
			}
			if ok {
				break
			}
			log.Info().Int64("ID", act.ID).Msg("querying activity details")
			act, err = input.Activity(ctx, act.ID)
			if err != nil {
				return err
			}
			n++
			log.Info().Int("n", n).Int64("ID", act.ID).Str("name", act.Name).Msg("saving activity details")
			if err = output.Save(ctx, act); err != nil {
				return err
			}
		}
	}
	return encoding.Encode(map[string]int{"total": total, "new": n, "existing": total - n})
}

func filterFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     "filter",
		Aliases:  []string{"f"},
		Required: required,
		Usage:    "Expression for filtering activities to remove",
	}
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
		filterFlag(true),
		&cli.BoolFlag{
			Name:    "dryrun",
			Aliases: []string{"n"},
			Value:   false,
			Usage:   "Don't actually remove anything, just show what would be done",
		},
	},
	Action: remove,
}

var exportCommand = &cli.Command{
	Name:   "export",
	Usage:  "Export activities from local storage",
	Flags:  []cli.Flag{filterFlag(false)},
	Action: export,
}

var Command = &cli.Command{
	Name:  "store",
	Usage: "Manage a local store of Strava activities",
	Flags: commands.Merge([]cli.Flag{commands.StoreFlag}, strava.AuthFlags),
	Subcommands: []*cli.Command{
		exportCommand,
		removeCommand,
		updateCommand,
	},
}
