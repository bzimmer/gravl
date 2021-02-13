package visualcrossing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bzimmer/gravl/pkg/providers/wx"
)

// ForecastService provides forecasts
type ForecastService service

func values(opts wx.ForecastOptions) (*url.Values, error) {
	v := &url.Values{}
	switch opts.AggregateHours {
	case 0:
		// do nothing
	case 1, 12, 24:
		v.Set("aggregateHours", strconv.FormatInt(int64(opts.AggregateHours), 10))
	default:
		return nil, &Fault{
			Message: fmt.Sprintf("unknown aggregate hours {%d}", opts.AggregateHours)}
	}
	switch {
	case opts.Location != "":
		v.Set("locations", opts.Location)
	case opts.Point != nil:
		loc := fmt.Sprintf("%0.4f,%0.4f", opts.Point.Y(), opts.Point.X())
		v.Set("locations", loc)
	default:
		return nil, &Fault{Message: "no location or coordinates specified"}
	}
	v.Set("includeAstronomy", fmt.Sprintf("%t", true))
	v.Set("alertLevel", "detail")
	switch opts.Units {
	case wx.Metric:
		v.Set("unitGroup", "metric")
	case wx.Imperial:
		v.Set("unitGroup", "us")
	}
	return v, nil
}

// Forecast weather conditions for a point
func (s *ForecastService) Forecast(ctx context.Context, opts wx.ForecastOptions) (*Forecast, error) {
	values, err := values(opts)
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
