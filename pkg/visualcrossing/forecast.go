package visualcrossing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	// AlertLevelNone .
	AlertLevelNone = "none"
	// AlertLevelSummary .
	AlertLevelSummary = "summary"
	// AlertLevelDetail .
	AlertLevelDetail = "detail"

	//UnitsUS .
	UnitsUS = "us"
	// UnitsUK .
	UnitsUK = "uk"
	// UnitsBase .
	UnitsBase = "base"
	// UnitsMetric .
	UnitsMetric = "metric"
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
		v.Add("aggregateHours", fmt.Sprintf("%d", hours))
		return nil
	}
}

// WithLocation .
func WithLocation(location string, locations ...string) ForecastOption {
	return func(v *url.Values) error {
		locs := []string{location}
		for _, n := range locations {
			locs = append(locs, n)
		}
		v.Set("locations", strings.Join(locs, "|"))
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
		v.Add("alertLevel", level)
		return nil
	}
}

// WithAstronomy .
func WithAstronomy(astro bool) ForecastOption {
	return func(v *url.Values) error {
		v.Add("includeAstronomy", fmt.Sprintf("%t", astro))
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
		case UnitsBase:
		default:
			return fmt.Errorf("unknown units {%s}", units)
		}
		v.Add("unitGroup", units)
		return nil
	}
}

// Forecast .
func (s *ForecastService) Forecast(ctx context.Context, opts ...ForecastOption) (*Forecast, error) {
	values, err := makeValues(opts)
	if err != nil {
		return nil, err
	}
	req, err := s.client.newAPIRequest(http.MethodGet, "forecast", values)
	if err != nil {
		return nil, err
	}
	fct := &Forecast{}
	err = s.client.Do(ctx, req, fct)
	if err != nil {
		return nil, err
	}
	return fct, err
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
