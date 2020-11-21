package openweather

import "github.com/bzimmer/gravl/pkg/common/wx"

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
	Day   float64 `json:"day"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type FeelsLike struct {
	Day   float64 `json:"day"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

type Current struct {
	Datetime   int        `json:"dt"`
	Sunrise    int        `json:"sunrise"`
	Sunset     int        `json:"sunset"`
	Temp       float64    `json:"temp"`
	FeelsLike  float64    `json:"feels_like"`
	Pressure   int        `json:"pressure"`
	Humidity   int        `json:"humidity"`
	DewPoint   float64    `json:"dew_point"`
	UVI        float64    `json:"uvi"`
	Clouds     int        `json:"clouds"`
	Visibility int        `json:"visibility"`
	WindSpeed  float64    `json:"wind_speed"`
	WindDeg    int        `json:"wind_deg"`
	WindGust   float64    `json:"wind_gust"`
	Weather    []*Weather `json:"weather"`
}

type Minutely struct {
	Datetime      int `json:"dt"`
	Precipitation int `json:"precipitation"`
}

type Hourly struct {
	Datetime   int        `json:"dt"`
	Temp       float64    `json:"temp"`
	FeelsLike  float64    `json:"feels_like"`
	Pressure   int        `json:"pressure"`
	Humidity   int        `json:"humidity"`
	DewPoint   float64    `json:"dew_point"`
	Clouds     int        `json:"clouds"`
	Visibility int        `json:"visibility"`
	WindSpeed  float64    `json:"wind_speed"`
	WindDeg    int        `json:"wind_deg"`
	Weather    []*Weather `json:"weather"`
	Pop        int        `json:"pop"`
	Rain       *Rain      `json:"rain,omitempty"`
}

type Rain struct {
	OneH float64 `json:"1h"`
}

type Daily struct {
	Datetime    int         `json:"dt"`
	Sunrise     int         `json:"sunrise"`
	Sunset      int         `json:"sunset"`
	Temperature Temperature `json:"temp"`
	FeelsLike   FeelsLike   `json:"feels_like"`
	Pressure    int         `json:"pressure"`
	Humidity    int         `json:"humidity"`
	DewPoint    float64     `json:"dew_point"`
	WindSpeed   float64     `json:"wind_speed"`
	WindDegree  int         `json:"wind_deg"`
	Weather     []*Weather  `json:"weather"`
	Clouds      int         `json:"clouds"`
	Pop         float64     `json:"pop"`
	Rain        float64     `json:"rain,omitempty"`
	UVI         float64     `json:"uvi"`
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

func (f *Forecast) WxForecast() (*wx.Forecast, error) {
	return nil, nil
}
