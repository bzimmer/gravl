package srtm

//go:generate go run ../../cmd/genwith/genwith.go --client --package srtm

import (
	"errors"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/tkrajina/go-elevations/geoelevations"
)

// Client client
type Client struct {
	storage geoelevations.SrtmLocalStorage
	client  *http.Client

	Elevation *ElevationService
}

func withServices(c *Client) {
	c.Elevation = &ElevationService{client: c}
}

// WithStorage sets the cache implementation for data files
func WithStorage(storage geoelevations.SrtmLocalStorage) Option {
	return func(c *Client) error {
		if storage == nil {
			return errors.New("nil storage")
		}
		c.storage = storage
		return nil
	}
}

// WithStorageLocation uses the fully qualified directory name to store cached files
func WithStorageLocation(directory string) Option {
	return func(c *Client) error {
		if directory == "" {
			return errors.New("nil directory")
		}
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			log.Info().Str("directory", directory).Msg("creating")
			if err := os.MkdirAll(directory, os.ModeDir|0700); err != nil {
				return err
			}
		}
		storage, err := geoelevations.NewLocalFileSrtmStorage(directory)
		if err != nil {
			return err
		}
		c.storage = storage
		log.Debug().Str("directory", directory).Msg("SRTM cache location")
		return nil
	}
}

func (c *Client) srtm() (*geoelevations.Srtm, error) {
	if c.storage == nil {
		return nil, errors.New("nil storage")
	}
	return geoelevations.NewSrtmWithCustomStorage(c.client, c.storage)
}
