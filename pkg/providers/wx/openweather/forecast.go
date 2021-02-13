package openweather

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bzimmer/gravl/pkg/providers/wx"
)

// ForecastService provides forecast
type ForecastService service

func values(opts wx.ForecastOptions) (*url.Values, error) {
	v := &url.Values{}
	if opts.Point == nil {
		return nil, &Fault{Message: "no coordinates specified"}
	}
	v.Set("lat", fmt.Sprintf("%0.4f", opts.Point.Y()))
	v.Set("lon", fmt.Sprintf("%0.4f", opts.Point.X()))
	switch opts.Units {
	case wx.Metric:
		v.Set("units", "metric")
	case wx.Imperial:
		v.Set("units", "imperial")
	}
	return v, nil
}

// ForecastOption .
type ForecastOption func(*url.Values) error

// Forecast returns a forecast
func (s *ForecastService) Forecast(ctx context.Context, opts wx.ForecastOptions) (*Forecast, error) {
	values, err := values(opts)
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
