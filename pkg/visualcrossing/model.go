package visualcrossing

import (
	"time"

	"github.com/bzimmer/gravl/pkg/common/wx"
)

// https://www.visualcrossing.com/resources/documentation/weather-api/weather-api-json-result-structure/

const (
	// AlertLevelNone .
	AlertLevelNone = "none"
	// AlertLevelSummary .
	AlertLevelSummary = "summary"
	// AlertLevelDetail .
	AlertLevelDetail = "detail"

	// UnitsUS .
	UnitsUS = "us"
	// UnitsUK .
	UnitsUK = "uk"
	// UnitsSI .
	UnitsSI = "base"
	// UnitsMetric .
	UnitsMetric = "metric"
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

type column struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
	Unit string `json:"unit"`
}

type columns struct {
	Address            column `json:"address"`
	CloudCover         column `json:"cloudcover"`
	Conditions         column `json:"conditions"`
	DateTime           column `json:"datetime"`
	Dew                column `json:"dew"`
	HeatIndex          column `json:"heatindex"`
	Humidity           column `json:"humidity"`
	Latitude           column `json:"latitude"`
	Longitude          column `json:"longitude"`
	LongwaveRadiation  column `json:"lw_radiation"`
	MaxTemp            column `json:"maxt"`
	MinTemp            column `json:"mint"`
	Name               column `json:"name"`
	PoP                column `json:"pop"`
	Precip             column `json:"precip"`
	ResolvedAddress    column `json:"resolvedAddress"`
	SeaLevelPressure   column `json:"sealevelpressure"`
	ShortwaveRadiation column `json:"sw_radiation"`
	Snow               column `json:"snow"`
	SnowDepth          column `json:"snowdepth"`
	Sunshine           column `json:"sunshine"`
	Temp               column `json:"temp"`
	Visibility         column `json:"visibility"`
	WindChill          column `json:"windchill"`
	WindDirection      column `json:"wdir"`
	WindGust           column `json:"wgust"`
	WindSpeed          column `json:"wspd"`
}

type alert struct {
	Event       string     `json:"event"`
	Headline    string     `json:"headline"`
	Description string     `json:"description"`
	Ends        *time.Time `json:"ends"`
	Onset       *time.Time `json:"onset"`
}

type conditions struct {
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

type forecastConditions struct {
	conditions
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
func (c *forecastConditions) WxConditions() *wx.Conditions {
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

type currentConditions struct {
	conditions
	DateTime  *time.Time `json:"datetime"`
	MoonPhase float64    `json:"moonphase"`
	Stations  string     `json:"stations"`
	Sunrise   *time.Time `json:"sunrise"`
	Sunset    *time.Time `json:"sunset"`
}

// WxConditions .
func (c *currentConditions) WxConditions() *wx.Conditions {
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

type location struct {
	Address            string               `json:"address"`
	Alerts             []alert              `json:"alerts"`
	ForecastConditions []forecastConditions `json:"values"`
	CurrentConditions  currentConditions    `json:"currentConditions"`
	Distance           float64              `json:"distance"`
	ID                 string               `json:"id"`
	Index              int                  `json:"index"`
	Latitude           float64              `json:"latitude"`
	Longitude          float64              `json:"longitude"`
	Name               string               `json:"name"`
	Time               float64              `json:"time"`
	Timezone           string               `json:"tz"`
}

type forecast struct {
	Columns       columns     `json:"columns"`
	RemainingCost int         `json:"remainingCost"`
	QueryCost     int         `json:"queryCost"`
	Messages      interface{} `json:"messages"`
	Locations     []location  `json:"locations"`
}

// WxForecast .
func (f *forecast) WxForecast() (*wx.Forecast, error) {
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
