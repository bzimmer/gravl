package openweather

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ForecastService .
type ForecastService service

// ForecastOption .
type ForecastOption func(*url.Values) error

// WithLocation .
func WithCoordinates(longitude, latitude float64) ForecastOption {
	return func(v *url.Values) error {
		v.Set("lat", fmt.Sprintf("%0.4f", latitude))
		v.Set("lon", fmt.Sprintf("%0.4f", longitude))
		return nil
	}
}

// WithUnits sets the unit family to use
func WithUnits(units Units) ForecastOption {
	return func(v *url.Values) error {
		var s string
		switch units {
		case UnitsImperial:
			s = "imperial"
		case UnitsMetric:
			s = "metric"
		case UnitsStandard:
			s = "standard"
		default:
			return fmt.Errorf("unknown units {%s}", units)
		}
		v.Set("units", s)
		return nil
	}
}

// Forecast returns a forecast
func (s *ForecastService) Forecast(ctx context.Context, opts ...ForecastOption) (*Forecast, error) {
	values, err := makeValues(opts)
	if err != nil {
		return nil, err
	}
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, "onecall", values)
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

func makeValues(options []ForecastOption) (*url.Values, error) {
	v := &url.Values{}
	for _, opt := range options {
		err := opt(v)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}
