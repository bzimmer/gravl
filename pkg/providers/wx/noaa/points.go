package noaa

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/twpayne/go-geom"
)

// PointsService .
type PointsService service

// float64ToString converts the float to the format necessary for NOAA
// It must have no more than four floating point values and no trailing zeros
func float64ToString(f float64) string {
	s := fmt.Sprintf("%0.2f", f)
	return strings.TrimRight(s, "0")
}

// GridPoint .
func (s *PointsService) GridPoint(ctx context.Context, point *geom.Point) (*GridPoint, error) {
	uri := fmt.Sprintf("points/%s,%s", float64ToString(point.Y()), float64ToString(point.X()))
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	gpt := &GridPoint{}
	err = s.client.do(req, gpt)
	if err != nil {
		return nil, err
	}
	return gpt, err
}

// Forecast .
func (s *PointsService) Forecast(ctx context.Context, point *geom.Point) (*Forecast, error) {
	uri := fmt.Sprintf("points/%s,%s/forecast", float64ToString(point.Y()), float64ToString(point.X()))
	return s.forecast(ctx, uri)
}

// ForecastHourly .
func (s *PointsService) ForecastHourly(ctx context.Context, point *geom.Point) (*Forecast, error) {
	uri := fmt.Sprintf("points/%s,%s/forecast/hourly", float64ToString(point.Y()), float64ToString(point.X()))
	return s.forecast(ctx, uri)
}

func (s *PointsService) forecast(ctx context.Context, uri string) (*Forecast, error) {
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, uri)
	if err != nil {
		return nil, err
	}
	fct := &Forecast{}
	err = s.client.do(req, fct)
	if err != nil {
		return nil, err
	}
	return fct, err
}