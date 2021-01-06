package srtm

import (
	"context"

	"github.com/twpayne/go-geom"
)

// ElevationService is the API for querying elevations from SRTM
type ElevationService service

// Elevation returns the elevation for the coordinates
func (s *ElevationService) Elevation(ctx context.Context, point *geom.Point) (float64, error) {
	m, err := s.client.srtm()
	if err != nil {
		return 0.0, err
	}
	return m.GetElevation(s.client.client, point.Y(), point.X())
}

// Elevations returns the elevations for the coordinates
func (s *ElevationService) Elevations(ctx context.Context, coords []*geom.Point) ([]float64, error) {
	m, err := s.client.srtm()
	if err != nil {
		return []float64{}, err
	}
	elevations := make([]float64, len(coords))
	for i, coord := range coords {
		elevation, err := m.GetElevation(s.client.client, coord.Y(), coord.X())
		if err != nil {
			return []float64{}, err
		}
		elevations[i] = elevation
	}
	return elevations, nil
}
