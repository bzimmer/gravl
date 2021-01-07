package strava

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
	"github.com/urfave/cli/v2"
	"go.etcd.io/bbolt"

	"github.com/bzimmer/gravl/pkg/activity/strava"
	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
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

	var x int64
	args := c.Args().Slice()
	ids := make([]interface{}, len(args))
	for i := 0; i < len(args); i++ {
		x, err = strconv.ParseInt(args[i], 0, 64)
		if err != nil {
			return err
		}
		ids[i] = x
	}

	q := bolthold.Where("ID").In(ids...)
	err = db.Bolt().Update(func(tx *bbolt.Tx) error {
		log.Info().Str("q", q.String()).Msg("deleting activities")
		if err = db.TxDeleteMatching(tx, strava.Activity{}, q); err != nil {
			return err
		}
		return nil
	})
	return err
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
	Name:  "update",
	Usage: "Query Strava for activities and update the local db",
	Flags: []cli.Flag{
		commands.StoreFlag,
		&cli.BoolFlag{
			Name:    "remove",
			Aliases: []string{"r"},
			Value:   false,
			Usage:   "Delete keys",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("remove") {
			return remove(c)
		}
		return update(c)
	},
}
