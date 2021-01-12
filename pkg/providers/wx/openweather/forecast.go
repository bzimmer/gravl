package openweather

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/twpayne/go-geom"
)

// ForecastService provides forecast
type ForecastService service

type ForecastOptions struct {
	Units Units
	Point *geom.Point
}

func (r *ForecastOptions) values() (*url.Values, error) {
	v := &url.Values{}
	if r.Point == nil {
		return nil, &Fault{Message: "no coordinates specified"}
	}
	v.Set("lat", fmt.Sprintf("%0.4f", r.Point.Y()))
	v.Set("lon", fmt.Sprintf("%0.4f", r.Point.X()))
	v.Set("units", r.Units.String())
	return v, nil
}

// ForecastOption .
type ForecastOption func(*url.Values) error

// Forecast returns a forecast
func (s *ForecastService) Forecast(ctx context.Context, opt ForecastOptions) (*Forecast, error) {
	values, err := opt.values()
	if err != nil {
		return nil, err
	}
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, "onecall", values)
	if err != nil {
		return nil, err
	}
	fct := &Forecast{}
	err = s.client.do(req, fct)
	if err != nil {
		return nil, err
	}
	return fct, nil
}
