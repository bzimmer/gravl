package strava

import (
	"github.com/twpayne/go-polyline"
)

func polylineToCoords(polylines ...string) ([][]float64, error) {
	zero := float64(0)
	var coords [][]float64
	for _, p := range polylines {
		if p == "" {
			continue
		}
		c, _, err := polyline.DecodeCoords([]byte(p))
		if err != nil {
			return nil, err
		}
		coords = make([][]float64, len(c))
		for i, x := range c {
			coords[i] = []float64{x[1], x[0], zero}
		}
		break
	}
	return coords, nil
}
