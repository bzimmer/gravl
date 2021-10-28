package pkg_test

import (
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-gpx"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/internal"
)

type something struct{}

func (s *something) GPX() (*gpx.GPX, error) {
	return &gpx.GPX{}, nil
}

func (s *something) GeoJSON() (*geojson.FeatureCollection, error) {
	return &geojson.FeatureCollection{}, nil
}

func TestRuntime(t *testing.T) {
	a := assert.New(t)
	tests := []*internal.Harness{
		{
			Name: "afters",
			Args: []string{"gravl", "afters"},
			After: pkg.Afters(
				func(c *cli.Context) error {
					pkg.Runtime(c).Metrics.IncrCounter([]string{"afters", "test"}, 1)
					pkg.Runtime(c).Metrics.AddSample([]string{"elapsed"}, float32(time.Since(pkg.Runtime(c).Start).Seconds()))
					return nil
				},
			),
			Counters: map[string]int{
				"gravl.afters.test": 1,
			},
		},
		{
			Name: "token",
			Args: []string{"gravl", "token"},
			Action: func(c *cli.Context) error {
				k, err := pkg.Token(16)
				if err != nil {
					return err
				}
				log.Info().Str("token", k).Msg(c.Command.Name)
				pkg.Runtime(c).Metrics.IncrCounter([]string{"token", "test"}, 1)
				return nil
			},
			Counters: map[string]int{
				"gravl.token.test": 1,
			},
		},
		{
			Name: "befores",
			Args: []string{"gravl", "befores"},
			Before: pkg.Befores(
				func(c *cli.Context) error {
					t, err := pkg.Token(16)
					if err != nil {
						return err
					}
					enc := pkg.JSON(c.App.Writer, false)
					if err := enc.Encode(t); err != nil {
						return err
					}
					enc = pkg.XML(c.App.Writer, false)
					if err := enc.Encode(t); err != nil {
						return err
					}
					pkg.Runtime(c).Metrics.IncrCounter([]string{"befores", "test"}, 1)
					return nil
				},
			),
			Counters: map[string]int{
				"gravl.befores.test": 1,
			},
		},
		{
			Name: "enc",
			Args: []string{"gravl", "enc"},
			Before: pkg.Befores(
				func(c *cli.Context) error {
					s := &something{}
					enc := pkg.JSON(c.App.Writer, false)
					if err := enc.Encode(s); err != nil {
						return err
					}
					enc = pkg.XML(c.App.Writer, false)
					if err := enc.Encode(s); err != nil {
						return err
					}
					enc = pkg.GPX(c.App.Writer, false)
					if err := enc.Encode(s); err != nil {
						return err
					}
					a.Error(enc.Encode(struct{}{}))
					enc = pkg.GeoJSON(c.App.Writer, false)
					if err := enc.Encode(s); err != nil {
						return err
					}
					a.Error(enc.Encode(struct{}{}))
					pkg.Runtime(c).Metrics.IncrCounter([]string{"enc", "test"}, 1)
					return nil
				},
			),
			Counters: map[string]int{
				"gravl.enc.test": 1,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			internal.Run(t, tt, nil, func(t *testing.T, baseURL string) *cli.Command {
				return &cli.Command{Name: tt.Name, Action: tt.Action}
			})
		})
	}
}
