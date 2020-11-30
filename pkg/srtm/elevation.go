package srtm

import "context"

// ElevationService is the API for querying elevations from SRTM
type ElevationService service

// Elevation returns the elevation for the coordinates
func (s *ElevationService) Elevation(ctx context.Context, longitude, latitude float64) (float64, error) {
	m, err := s.client.srtm()
	if err != nil {
		return 0.0, err
	}
	return m.GetElevation(s.client.client, latitude, longitude)
}

// Elevations returns the elevations for the coordinates
func (s *ElevationService) Elevations(ctx context.Context, coords [][]float64) ([]float64, error) {
	m, err := s.client.srtm()
	if err != nil {
		return []float64{}, err
	}
	elevations := make([]float64, len(coords))
	for i, coord := range coords {
		// input is expected in (longitude, latitude) format so swap
		elevation, err := m.GetElevation(s.client.client, coord[1], coord[0])
		if err != nil {
			return []float64{}, err
		}
		elevations[i] = elevation
	}
	return elevations, nil
}
