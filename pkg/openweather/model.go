package openweather

//go:generate stringer -type=Units -linecomment -output=model_string.go

import (
	"errors"

	"github.com/bzimmer/gravl/pkg/common/wx"
)

// Units of measure
//  https://openweathermap.org/api/one-call-api#data
type Units int

const (
	// UnitsMetric fo temperature in Celsius and wind speed in meter/sec
	UnitsMetric Units = iota // metric
	// UnitsImperial for temperature in Fahrenheit and wind speed in miles/hour
	UnitsImperial // imperial
	// UnitsStandard for temperature in Kelvin and wind speed in meter/sec
	UnitsStandard // standard
)

// Fault .
type Fault struct {
	ErrorCode int    `json:"cod"`
	Message   string `json:"message"`
}

func (f *Fault) Error() string {
	return f.Message
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Temperature struct {
	Day     float64 `json:"day"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Night   float64 `json:"night"`
	Evening float64 `json:"eve"`
	Morning float64 `json:"morn"`
}

type FeelsLike struct {
	Day     float64 `json:"day"`
	Night   float64 `json:"night"`
	Evening float64 `json:"eve"`
	Morning float64 `json:"morn"`
}

type Current struct {
	Datetime    int        `json:"dt"`
	Sunrise     int        `json:"sunrise"`
	Sunset      int        `json:"sunset"`
	Temperature float64    `json:"temp"`
	FeelsLike   float64    `json:"feels_like"`
	Pressure    float64    `json:"pressure"`
	Humidity    float64    `json:"humidity"`
	DewPoint    float64    `json:"dew_point"`
	UVI         float64    `json:"uvi"`
	Clouds      float64    `json:"clouds"`
	Visibility  float64    `json:"visibility"`
	WindSpeed   float64    `json:"wind_speed"`
	WindDeg     float64    `json:"wind_deg"`
	WindGust    float64    `json:"wind_gust"`
	Weather     []*Weather `json:"weather"`
}

type Minutely struct {
	Datetime      int     `json:"dt"`
	Precipitation float64 `json:"precipitation"`
}

type Hourly struct {
	Datetime          int        `json:"dt"`
	Temperature       float64    `json:"temp"`
	FeelsLike         float64    `json:"feels_like"`
	Pressure          float64    `json:"pressure"`
	Humidity          float64    `json:"humidity"`
	DewPoint          float64    `json:"dew_point"`
	Clouds            float64    `json:"clouds"`
	Visibility        float64    `json:"visibility"`
	WindSpeed         float64    `json:"wind_speed"`
	WindDegree        float64    `json:"wind_deg"`
	Weather           []*Weather `json:"weather"`
	PrecipProbability float64    `json:"pop"`
	Rain              *Rain      `json:"rain,omitempty"`
}

type Rain struct {
	OneH float64 `json:"1h"`
}

type Daily struct {
	Datetime          int         `json:"dt"`
	Sunrise           int         `json:"sunrise"`
	Sunset            int         `json:"sunset"`
	Temperature       Temperature `json:"temp"`
	FeelsLike         FeelsLike   `json:"feels_like"`
	Pressure          float64     `json:"pressure"`
	Humidity          float64     `json:"humidity"`
	DewPoint          float64     `json:"dew_point"`
	WindSpeed         float64     `json:"wind_speed"`
	WindDegree        float64     `json:"wind_deg"`
	Weather           []*Weather  `json:"weather"`
	Clouds            float64     `json:"clouds"`
	PrecipProbability float64     `json:"pop"`
	Rain              float64     `json:"rain,omitempty"`
	UVI               float64     `json:"uvi"`
}

type Forecast struct {
	Latitude       float64     `json:"lat"`
	Longitude      float64     `json:"lon"`
	Timezone       string      `json:"timezone"`
	TimezoneOffset int         `json:"timezone_offset"`
	Current        *Current    `json:"current"`
	Minutely       []*Minutely `json:"minutely"`
	Hourly         []*Hourly   `json:"hourly"`
	Daily          []*Daily    `json:"daily"`
}

func (f *Forecast) Forecast() (*wx.Forecast, error) {
	return nil, errors.New("not implemented")
}
