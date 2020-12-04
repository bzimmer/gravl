package visualcrossing

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
	AggregateHours int
	Location       string
	Coordinates    Coordinates
	AlertLevel     AlertLevel
	Astronomy      bool
	Units          Units
}

func (r *ForecastOptions) values() (*url.Values, error) {
	v := &url.Values{}
	switch r.AggregateHours {
	case 0:
		// do nothing
	case 1, 12, 24:
		v.Set("aggregateHours", fmt.Sprintf("%d", r.AggregateHours))
	default:
		return nil, &Fault{
			Message: fmt.Sprintf("unknown aggregate hours {%d}", r.AggregateHours)}
	}
	switch {
	case r.Location != "":
		v.Set("locations", r.Location)
	case r.Coordinates.Latitude != 0.0 && r.Coordinates.Longitude != 0.0:
		v.Set("locations", r.Coordinates.String())
	default:
		return nil, &Fault{Message: "no location or coordinates specified"}
	}
	v.Set("includeAstronomy", fmt.Sprintf("%t", r.Astronomy))
	v.Set("alertLevel", r.AlertLevel.String())
	v.Set("unitGroup", r.Units.String())
	return v, nil
}

// Forecast .
func (s *ForecastService) Forecast(ctx context.Context, opt ForecastOptions) (*Forecast, error) {
	values, err := opt.values()
	if err != nil {
		return nil, err
	}
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, "forecast", values)
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
