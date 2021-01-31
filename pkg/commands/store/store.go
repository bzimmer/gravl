package store

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	stravaapi "github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func filter(c *cli.Context, acts []*stravaapi.Activity) ([]*stravaapi.Activity, error) {
	if !c.IsSet("filter") {
		return acts, nil
	}
	evaluator := commands.Filterer(c.String("filter"))
	return evaluator.Filter(c.Context, acts)
}

func export(c *cli.Context) error {
	db, err := Open(c, "input")
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
	db, err := Open(c, "input")
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
	in, err := Open(c, "input")
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := Open(c, "output")
	if err != nil {
		return err
	}
	defer out.Close()
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()
	acts, errs := in.Activities(ctx)
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
			ok, err = out.Exists(ctx, act.ID)
			if err != nil {
				return err
			}
			if ok {
				break
			}
			log.Info().Int64("ID", act.ID).Msg("querying activity details")
			act, err = in.Activity(ctx, act.ID)
			if err != nil {
				return err
			}
			n++
			log.Info().Int("n", n).Int64("ID", act.ID).Str("name", act.Name).Msg("saving activity details")
			if err = out.Save(ctx, act); err != nil {
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
	Flags:  append([]cli.Flag{InputFlag("strava"), OutputFlag(DefaultLocalStore)}, strava.AuthFlags...),
}

var removeCommand = &cli.Command{
	Name:  "remove",
	Usage: "Remove activities from local storage",
	Flags: []cli.Flag{
		InputFlag(DefaultLocalStore),
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
	Flags:  []cli.Flag{InputFlag(DefaultLocalStore), filterFlag(false)},
	Action: export,
}

func InputFlag(storeDefault string) cli.Flag {
	return &cli.StringFlag{
		Name:    "input",
		Aliases: []string{"i"},
		Value:   storeDefault,
		Usage:   "Input data store"}
}

func OutputFlag(storeDefault string) cli.Flag {
	return &cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Value:   storeDefault,
		Usage:   "Output data store",
	}
}

var Command = &cli.Command{
	Name:  "store",
	Usage: "Manage a local store of Strava activities",
	Subcommands: []*cli.Command{
		exportCommand,
		removeCommand,
		updateCommand,
	},
}
