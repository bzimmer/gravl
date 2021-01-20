package staticmap

import (
	"context"
	"flag"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"sync"

	"github.com/adrg/xdg"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/rs/zerolog/log"
	"github.com/twpayne/go-geom"
	"golang.org/x/sync/errgroup"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/analysis"
	"github.com/bzimmer/gravl/pkg/providers/activity/strava"
)

const doc = `staticmap generates a staticmap for every activity.`

type smap struct {
	output  string
	workers int
}

type activity struct {
	a *strava.Activity
	s *geom.LineString
	p string
}

func (s *smap) linestrings(ctx context.Context, g *errgroup.Group, acts []*strava.Activity) <-chan *activity {
	activities := make(chan *activity)
	g.Go(func() error {
		defer close(activities)
		for i := range acts {
			act := acts[i]
			if act.Map == nil {
				continue
			}
			trk, err := act.Map.LineString()
			if err != nil {
				log.Error().Str("name", act.Name).Err(err).Msg("linestring")
				continue
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case activities <- &activity{a: act, s: trk}:
			}
		}
		return nil
	})
	return activities
}

func (s *smap) paths(ctx context.Context, g *errgroup.Group, activities <-chan *activity) <-chan *activity {
	var wg sync.WaitGroup
	paths := make(chan *activity)

	defer func() {
		go func() {
			wg.Wait()
			close(paths)
		}()
	}()

	for i := 0; i < s.workers; i++ {
		g.Go(func() error {
			wg.Add(1)
			defer func() { wg.Done() }()
			for act := range activities {
				log.Info().Str("name", act.a.Name).Msg("creating staticmap")
				ictx := sm.NewContext()
				ictx.SetSize(1280, 800)
				for _, coord := range act.s.Coords() {
					pt := s2.LatLngFromDegrees(coord.Y(), coord.X())
					ictx.AddMarker(sm.NewMarker(pt, color.RGBA{0xff, 0, 0, 0xff}, 1.0))
				}
				img, err := ictx.Render()
				if err != nil {
					return err
				}
				act.p = filepath.Join(s.output, fmt.Sprintf("%d.png", act.a.ID))
				if err := gg.SavePNG(act.p, img); err != nil {
					return err
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case paths <- act:
				}
			}
			return nil
		})
	}
	return paths
}

func (s *smap) run(ctx *analysis.Context, pass *analysis.Pass) (interface{}, error) {
	s.output = filepath.Join(xdg.CacheHome, pkg.PackageName, "passes", "staticmap")
	if err := os.MkdirAll(s.output, os.ModeDir|0700); err != nil {
		return nil, err
	}

	g, c := errgroup.WithContext(ctx.Context)
	paths := s.paths(c, g, s.linestrings(c, g, pass.Activities))

	res := make(map[int64]string)
	g.Go(func() error {
		for {
			select {
			case <-c.Done():
				return c.Err()
			case act := <-paths:
				if act == nil {
					return nil
				}
				res[act.a.ID] = act.p
			}
		}
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return res, nil
}

func New() *analysis.Analyzer {
	// @todo(bzimmer) add flags
	s := &smap{
		workers: 15,
	}
	fs := flag.NewFlagSet("staticmap", flag.ExitOnError)
	fs.IntVar(&s.workers, "workers", s.workers, "number of workers")
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   doc,
		Flags: fs,
		Run:   s.run,
	}
}
