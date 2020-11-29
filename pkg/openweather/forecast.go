package openweather

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ForecastService .
type ForecastService service

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

func (c Coordinates) String() string {
	return fmt.Sprintf("%0.4f,%0.4f", c.Latitude, c.Longitude)
}

type ForecastOptions struct {
	Coordinates Coordinates
	Units       Units
}

func (r *ForecastOptions) values() (*url.Values, error) {
	v := &url.Values{}
	if r.Coordinates.Latitude == 0.0 && r.Coordinates.Longitude == 0.0 {
		return nil, &Fault{Message: "no coordinates specified"}
	}
	v.Set("locations", r.Coordinates.String())
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
	err = s.client.Do(req, fct)
	if err != nil {
		return nil, err
	}
	return fct, nil
}
