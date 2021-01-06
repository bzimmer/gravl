package strava

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/bzimmer/gravl/pkg/analysis/store"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
	"github.com/urfave/cli/v2"
)

func update(c *cli.Context) error {
	client, err := NewAPIClient(c)
	if err != nil {
		return err
	}
	fn := c.Path("db")
	if fn == "" {
		return errors.New("nil db path")
	}
	directory := filepath.Dir(fn)
	if _, err = os.Stat(directory); os.IsNotExist(err) {
		log.Info().Str("directory", directory).Msg("creating")
		if err = os.MkdirAll(directory, os.ModeDir|0700); err != nil {
			return err
		}
	}
	db, err := bolthold.Open(fn, 0666, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Info().Str("db", fn).Msg("using database")

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
	Usage:  "Query Strava for activities and update the local db",
	Action: update,
}
