package strava

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
	"github.com/urfave/cli/v2"
	"go.etcd.io/bbolt"

	"github.com/bzimmer/gravl/pkg/analysis/eval/antonmedv"
	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

func database(c *cli.Context) (*bolthold.Store, error) {
	fn := c.Path("store")
	if fn == "" {
		return nil, errors.New("nil db path")
	}
	log.Info().Str("store", fn).Msg("using database")
	directory := filepath.Dir(fn)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		log.Info().Str("directory", directory).Msg("creating")
		if err = os.MkdirAll(directory, os.ModeDir|0700); err != nil {
			return nil, err
		}
	}
	return bolthold.Open(fn, 0666, nil)
}

func remove(c *cli.Context) error {
	db, err := database(c)
	if err != nil {
		return err
	}
	defer db.Close()

	var acts []*strava.Activity
	err = db.ForEach(&bolthold.Query{}, func(act *strava.Activity) error {
		acts = append(acts, act)
		return nil
	})
	if c.IsSet("filter") {
		evaluator := antonmedv.New()
		acts, err = evaluator.Filter(c.Context, c.String("filter"), acts)
		if err != nil {
			return err
		}
	}
	ids := make([]interface{}, len(acts))
	for i, act := range acts {
		ids[i] = act.ID
	}
	if c.Bool("dryrun") {
		return encoding.Encode(ids)
	}
	q := bolthold.Where("ID").In(ids...)
	log.Info().Str("q", q.String()).Msg("deleting activities")
	err = db.Bolt().Update(func(tx *bbolt.Tx) error {
		if err = db.TxDeleteMatching(tx, strava.Activity{}, q); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return encoding.Encode(ids)
}

func update(c *cli.Context) error {
	client, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	db, err := database(c)
	if err != nil {
		return err
	}
	defer db.Close()

	var source store.Source
	if c.NArg() == 1 {
		source = &store.SourceFile{Path: c.Args().First()}
	} else {
		source = &store.SourceStrava{Client: client}
	}

	store := store.NewStore(db)
	n, err := store.Update(c.Context, source)
	if err != nil {
		return err
	}
	if err = encoding.Encode(map[string]int{"activities": n}); err != nil {
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
