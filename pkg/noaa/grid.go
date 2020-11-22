package noaa

import (
	"context"
	"fmt"
	"net/http"
)

// GridPointsService .
type GridPointsService service

// Forecast .
func (s *GridPointsService) Forecast(ctx context.Context, wfo string, x, y int) (*Forecast, error) {
	uri := fmt.Sprintf("gridpoints/%s/%d,%d/forecast", wfo, x, y)
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	fct := &Forecast{}
	err = s.client.Do(req, fct)
	if err != nil {
		return nil, err
	}
	return fct, nil
}
