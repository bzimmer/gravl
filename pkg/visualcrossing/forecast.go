package visualcrossing

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

// WithAggregateHours .
func WithAggregateHours(hours int) ForecastOption {
	return func(v *url.Values) error {
		switch hours {
		case 1:
		case 12:
		case 24:
		default:
			return fmt.Errorf("unknown aggregage hours {%d}", hours)
		}
		v.Set("aggregateHours", fmt.Sprintf("%d", hours))
		return nil
	}
}

// WithLocation .
func WithLocation(location string) ForecastOption {
	// VC supports more than one location per call but we're going
	//  to limit for sake of simplicity right now
	return func(v *url.Values) error {
		v.Set("locations", location)
		return nil
	}
}

// WithAlerts .
func WithAlerts(level string) ForecastOption {
	return func(v *url.Values) error {
		switch level {
		case AlertLevelNone:
		case AlertLevelSummary:
		case AlertLevelDetail:
		default:
			return fmt.Errorf("unknown alert level {%s}", level)
		}
		v.Set("alertLevel", level)
		return nil
	}
}

// WithAstronomy .
func WithAstronomy(astro bool) ForecastOption {
	return func(v *url.Values) error {
		v.Set("includeAstronomy", fmt.Sprintf("%t", astro))
		return nil
	}
}

// WithUnits .
func WithUnits(units string) ForecastOption {
	return func(v *url.Values) error {
		switch units {
		case UnitsUS:
		case UnitsUK:
		case UnitsMetric:
		case UnitsSI:
		default:
			return fmt.Errorf("unknown units {%s}", units)
		}
		v.Set("unitGroup", units)
		return nil
	}
}

// Forecast .
func (s *ForecastService) Forecast(ctx context.Context, opts ...ForecastOption) (*Forecast, error) {
	values, err := makeValues(opts)
	if err != nil {
		return nil, err
	}
	req, err := s.client.newAPIRequest(ctx, http.MethodGet, "forecast", values)
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
