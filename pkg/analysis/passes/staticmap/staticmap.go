package staticmap

import (
	"context"
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/rs/zerolog/log"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/analysis"
)

const Doc = `staticmap generates a staticmap for every activity`

func Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	paths := make(map[int64]string)
	output := filepath.Join(xdg.CacheHome, pkg.PackageName, "passes", "staticmap")
	if err := os.MkdirAll(output, os.ModeDir|0700); err != nil {
		return nil, err
	}
	for i := range pass.Activities {
		ctx := sm.NewContext()
		ctx.SetSize(1280, 800)
		act := pass.Activities[i]
		if act.Map == nil {
			continue
		}
		trk, err := act.Map.LineString()
		if err != nil {
			return nil, err
		}
		log.Info().Str("name", act.Name).Msg("creating staticmap")
		for _, coord := range trk.Coords() {
			pt := s2.LatLngFromDegrees(coord.Y(), coord.X())
			ctx.AddMarker(sm.NewMarker(pt, color.RGBA{0xff, 0, 0, 0xff}, 1.0))
		}
		img, err := ctx.Render()
		if err != nil {
			return nil, err
		}
		path := filepath.Join(output, fmt.Sprintf("%d.png", act.ID))
		if err := gg.SavePNG(path, img); err != nil {
			return nil, err
		}
		paths[act.ID] = path
	}
	return paths, nil
}

func New() *analysis.Analyzer {
	// @todo(bzimmer) add flags
	return &analysis.Analyzer{
		Name: "staticmap",
		Doc:  Doc,
		Run:  Run,
	}
}
