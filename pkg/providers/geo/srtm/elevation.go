package srtm

import (
	"context"

	"github.com/tkrajina/go-elevations/geoelevations"
	"github.com/twpayne/go-geom"
)

// ElevationService is the API for querying elevations from SRTM
type ElevationService struct {
	service
	srtm *geoelevations.Srtm
}

// Elevation returns the elevation for the coordinates
func (s *ElevationService) Elevation(ctx context.Context, point *geom.Point) (float64, error) {
	return s.srtm.GetElevation(s.client.client, point.Y(), point.X())
}

// Elevations returns the elevations for the coordinates
func (s *ElevationService) Elevations(ctx context.Context, coords []*geom.Point) ([]float64, error) {
	elevations := make([]float64, len(coords))
	for i, coord := range coords {
		elevation, err := s.srtm.GetElevation(s.client.client, coord.Y(), coord.X())
		if err != nil {
			return []float64{}, err
		}
		elevations[i] = elevation
	}
	return elevations, nil
}
