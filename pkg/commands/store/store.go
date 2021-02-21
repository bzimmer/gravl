package store

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/activity/strava"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/eval"
	stravaapi "github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func evaluator(c *cli.Context) (eval.Evaluator, error) {
	if !c.IsSet("filter") {
		return nil, nil
	}
	return commands.Evaluator(c.String("filter"))
}

func filter(ctx context.Context, filterer eval.Evaluator, acts <-chan *stravaapi.ActivityResult) <-chan *stravaapi.ActivityResult {
	res := make(chan *stravaapi.ActivityResult, 1)
	go func() {
		defer close(res)
		for {
			select {
			case <-ctx.Done():
				res <- &stravaapi.ActivityResult{Err: ctx.Err()}
				return
			case x, ok := <-acts:
				if !ok {
					return
				}
				if x.Err != nil {
					res <- &stravaapi.ActivityResult{Err: x.Err}
					continue
				}
				if filterer == nil {
					res <- &stravaapi.ActivityResult{Activity: x.Activity}
					continue
				}
				b, err := filterer.Bool(ctx, x.Activity)
				if err != nil {
					res <- &stravaapi.ActivityResult{Err: x.Err}
				} else if b {
					res <- &stravaapi.ActivityResult{Activity: x.Activity}
				}
			}
		}
	}()
	return res
}

func attributer(c *cli.Context) (func(ctx context.Context, act *stravaapi.Activity) (interface{}, error), error) {
	f := func(ctx context.Context, act *stravaapi.Activity) (interface{}, error) { return act, nil }
	if c.IsSet("attribute") {
		var evaluator eval.Evaluator
		evaluator, err := commands.Evaluator(c.String("attribute"))
		if err != nil {
			return nil, err
		}
		f = evaluator.Eval
	}
	return f, nil
}

func export(c *cli.Context) error {
	evaluator, err := evaluator(c)
	if err != nil {
		return err
	}
	db, err := Open(c, "input")
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()
	acts := db.Activities(ctx)
	attr, err := attributer(c)
	if err != nil {
		return err
	}
	acts = filter(ctx, evaluator, acts)
	if err != nil {
		return err
	}
	var i int
	defer func(t time.Time) {
		log.Info().Int("activities", i).Dur("elapsed", time.Since(t)).Msg("export")
	}(time.Now())
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case x, ok := <-acts:
			if !ok {
				return nil
			}
			if x.Err != nil {
				return x.Err
			}
			y, err := attr(ctx, x.Activity)
			if err != nil {
				return err
			}
			if err := encoding.Encode(y); err != nil {
				return err
			}
			i++
		}
	}
}

func remove(c *cli.Context) error {
	evaluator, err := evaluator(c)
	if err != nil {
		return err
	}
	db, err := Open(c, "input")
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()
	acts := db.Activities(ctx)
	acts = filter(ctx, evaluator, acts)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case x, ok := <-acts:
			if !ok {
				return nil
			}
			if x.Err != nil {
				return x.Err
			}
			if err := encoding.Encode(x.Activity.ID); err != nil {
				return err
			}
			if !c.Bool("dryrun") {
				if err := db.Remove(c.Context, x.Activity); err != nil {
					return err
				}
			}
		}
	}
}

func update(c *cli.Context) error {
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
	acts := in.Activities(ctx)
	for active := true; active; {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res, ok := <-acts:
			if !ok {
				// break the loop to return the processing results
				active = false
				break
			}
			if res.Err != nil {
				return res.Err
			}
			total++
			ok, err = out.Exists(ctx, res.Activity.ID)
			if err != nil {
				return err
			}
			if ok {
				break
			}
			log.Info().Int64("ID", res.Activity.ID).Msg("querying activity details")
			act, err := in.Activity(ctx, res.Activity.ID)
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
	Name:  "export",
	Usage: "Export activities from local storage",
	Flags: []cli.Flag{
		InputFlag(DefaultLocalStore),
		filterFlag(false),
		&cli.StringSliceFlag{
			Name:    "attribute",
			Aliases: []string{"B"},
			Usage:   "Evaluate the expression on an activity and return only those results",
		},
	},
	Action: export,
}

func filterFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     "filter",
		Aliases:  []string{"f"},
		Required: required,
		Usage:    "Expression for filtering activities",
	}
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
