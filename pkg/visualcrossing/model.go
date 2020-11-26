package visualcrossing

//go:generate stringer -type=Units,AlertLevel -output model_string.go

import (
	"time"

	"github.com/bzimmer/gravl/pkg/common/wx"
)

// https://www.visualcrossing.com/resources/documentation/weather-api/weather-api-json-result-structure/

// Units of measure
//  https://www.visualcrossing.com/resources/documentation/weather-api/unit-groups-and-measurement-units/
type Units int

// AlertLevel of alerts, warnings and other high priority information issued by local weather organizations
//  https://www.visualcrossing.com/resources/documentation/weather-data/weather-alerts/
type AlertLevel int

const (
	// UnitsUS for temperature in Fahrenheit and wind speed in miles/hour
	UnitsUS Units = iota
	// UnitsUK for temperature in Celsius and wind speed in miles/hour
	UnitsUK
	// UnitsStandard for temperature in Kelvin and wind speed in meter/sec
	UnitsStandard
	// UnitsMetric for temperature in Celsius and wind speed in km/hour
	UnitsMetric

	// None does not retrieve alert information (equivalent of omitting the parameter)
	AlertLevelNone AlertLevel = iota
	// Summary does not retrieve the detail field text of the alert
	AlertLevelSummary
	// Detail returns a full description of the alert including the detail
	AlertLevelDetail
)

// Fault .
type Fault struct {
	ErrorCode     int    `json:"errorCode"`
	ExecutionTime int    `json:"executionTime"`
	Message       string `json:"message"`
	SessionID     string `json:"sessionId"`
}

func (f *Fault) Error() string {
	return f.Message
}

type Column struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
	Unit string `json:"unit"`
}

type Columns struct {
	Address            Column `json:"address"`
	CloudCover         Column `json:"cloudcover"`
	Conditions         Column `json:"conditions"`
	DateTime           Column `json:"datetime"`
	Dew                Column `json:"dew"`
	HeatIndex          Column `json:"heatindex"`
	Humidity           Column `json:"humidity"`
	Latitude           Column `json:"latitude"`
	Longitude          Column `json:"longitude"`
	LongwaveRadiation  Column `json:"lw_radiation"`
	MaxTemp            Column `json:"maxt"`
	MinTemp            Column `json:"mint"`
	Name               Column `json:"name"`
	PoP                Column `json:"pop"`
	Precip             Column `json:"precip"`
	ResolvedAddress    Column `json:"resolvedAddress"`
	SeaLevelPressure   Column `json:"sealevelpressure"`
	ShortwaveRadiation Column `json:"sw_radiation"`
	Snow               Column `json:"snow"`
	SnowDepth          Column `json:"snowdepth"`
	Sunshine           Column `json:"sunshine"`
	Temp               Column `json:"temp"`
	Visibility         Column `json:"visibility"`
	WindChill          Column `json:"windchill"`
	WindDirection      Column `json:"wdir"`
	WindGust           Column `json:"wgust"`
	WindSpeed          Column `json:"wspd"`
}

type Alert struct {
	Event       string     `json:"event"`
	Headline    string     `json:"headline"`
	Description string     `json:"description"`
	Ends        *time.Time `json:"ends"`
	Onset       *time.Time `json:"onset"`
}

type Conditions struct {
	CloudCover       float64 `json:"cloudcover"`
	Dew              float64 `json:"dew"`
	HeatIndex        float64 `json:"heatindex"`
	Humidity         float64 `json:"humidity"`
	Precip           float64 `json:"precip"`
	SeaLevelPressure float64 `json:"sealevelpressure"`
	SnowDepth        float64 `json:"snowdepth"`
	Temperature      float64 `json:"temp"`
	Visibility       float64 `json:"visibility"`
	WindChill        float64 `json:"windchill"`
	WindBearing      float64 `json:"wdir"`
	WindGust         float64 `json:"wgust"`
	WindSpeed        float64 `json:"wspd"`
}

type ForecastConditions struct {
	Conditions
	Description        string     `json:"conditions"`
	DateTime           *time.Time `json:"datetimeStr"`
	LongWaveRadiation  float64    `json:"lw_radiation"`
	MaxTemp            float64    `json:"maxt"`
	MinTemp            float64    `json:"mint"`
	PrecipProbability  float64    `json:"pop"`
	Snow               float64    `json:"snow"`
	Sunshine           float64    `json:"sunshine"`
	ShortWaveRadiation float64    `json:"sw_radiation"`
}

// WxConditions .
func (c *ForecastConditions) WxConditions() *wx.Conditions {
	return &wx.Conditions{
		ValidFrom:         c.DateTime,
		Summary:           c.Description,
		Temperature:       c.Temperature,
		WindBearing:       c.WindBearing,
		WindChill:         c.WindChill,
		WindGust:          c.WindGust,
		WindSpeed:         c.WindSpeed,
		Precip:            c.Precip,
		PrecipProbability: c.PrecipProbability,
		TemperatureMax:    c.MaxTemp,
		TemperatureMin:    c.MinTemp,
	}
}

type CurrentConditions struct {
	Conditions
	DateTime  *time.Time `json:"datetime"`
	MoonPhase float64    `json:"moonphase"`
	Stations  string     `json:"stations"`
	Sunrise   *time.Time `json:"sunrise"`
	Sunset    *time.Time `json:"sunset"`
}

// WxConditions .
func (c *CurrentConditions) WxConditions() *wx.Conditions {
	return &wx.Conditions{
		ValidFrom:   c.DateTime,
		Sunrise:     c.Sunrise,
		Sunset:      c.Sunset,
		Temperature: c.Temperature,
		WindBearing: c.WindBearing,
		WindChill:   c.WindChill,
		WindGust:    c.WindGust,
		WindSpeed:   c.WindSpeed,
		Precip:      c.Precip,
		MoonPhase:   c.MoonPhase,
		DewPoint:    c.Dew,
		CloudCover:  c.CloudCover,
		SnowDepth:   c.SnowDepth,
	}
}

type Location struct {
	Address            string                `json:"address"`
	Alerts             []*Alert              `json:"alerts"`
	ForecastConditions []*ForecastConditions `json:"values"`
	CurrentConditions  *CurrentConditions    `json:"currentConditions"`
	Distance           float64               `json:"distance"`
	ID                 string                `json:"id"`
	Index              int                   `json:"index"`
	Latitude           float64               `json:"latitude"`
	Longitude          float64               `json:"longitude"`
	Name               string                `json:"name"`
	Time               float64               `json:"time"`
	Timezone           string                `json:"tz"`
}

type Forecast struct {
	Columns       Columns     `json:"columns"`
	RemainingCost int         `json:"remainingCost"`
	QueryCost     int         `json:"queryCost"`
	Messages      interface{} `json:"messages"`
	Locations     []Location  `json:"locations"`
}

// Forecast .
func (f *Forecast) Forecast() (*wx.Forecast, error) {
	loc := f.Locations[0]
	conditions := make([]*wx.Conditions, len(loc.ForecastConditions))
	for i, c := range loc.ForecastConditions {
		conditions[i] = c.WxConditions()
	}
	return &wx.Forecast{
		ID:        loc.ID,
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
		Timezone:  loc.Timezone,
		Current:   loc.CurrentConditions.WxConditions(),
		Period:    &wx.Period{Conditions: conditions},
	}, nil
}
