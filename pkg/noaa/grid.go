package noaa

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bzimmer/gravl/pkg/common/wx"
)

// GridPointsService .
type GridPointsService service

// Forecast .
func (s *GridPointsService) Forecast(ctx context.Context, wfo string, x, y int) (*wx.Forecast, error) {
	uri := fmt.Sprintf("gridpoints/%s/%d,%d/forecast", wfo, x, y)
	return s.forecast(ctx, uri)
}

func (s *GridPointsService) forecast(ctx context.Context, uri string) (*wx.Forecast, error) {
	req, err := s.client.newAPIRequest(http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	fct := &forecast{}
	err = s.client.Do(ctx, req, fct)
	if err != nil {
		return nil, err
	}
	return fct.WxForecasts()
}
