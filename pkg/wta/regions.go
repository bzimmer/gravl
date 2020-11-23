package wta

import "context"

//go:generate go run genregions.go

// RegionsService .
type RegionsService service

// Regions .
func (s *RegionsService) Regions(ctx context.Context) ([]*Region, error) {
	return regions, nil
}
