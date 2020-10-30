package wx

import (
	"time"

	gj "github.com/paulmach/go.geojson"
)

// Units .
type Units string

const (
	// SI .
	SI Units = "si"
	// US .
	US Units = "us"
	// UK .
	UK Units = "uk"
)

// Conditions .
type Conditions struct {
	ValidFrom                  *time.Time `json:"validFrom,omitempty"`
	ValidTo                    *time.Time `json:"validTo,omitempty"`
	Summary                    string     `json:"summary,omitempty"`
	Icon                       string     `json:"icon,omitempty"`
	Sunrise                    *time.Time `json:"sunriseTime,omitempty"`
	Sunset                     *time.Time `json:"sunsetTime,omitempty"`
	Precip                     float64    `json:"precip,omitempty"`
	PrecipIntensity            float64    `json:"precipIntensity,omitempty"`
	PrecipIntensityMax         float64    `json:"precipIntensityMax,omitempty"`
	PrecipIntensityMaxTime     *time.Time `json:"precipIntensityMaxTime,omitempty"`
	PrecipProbability          float64    `json:"precipProbability,omitempty"`
	PrecipType                 string     `json:"precipType,omitempty"`
	PrecipAccumulation         float64    `json:"precipAccumulation,omitempty"`
	Temperature                float64    `json:"temperature,omitempty"`
	TemperatureMin             float64    `json:"temperatureMin,omitempty"`
	TemperatureMinTime         *time.Time `json:"temperatureMinTime,omitempty"`
	TemperatureMax             float64    `json:"temperatureMax,omitempty"`
	TemperatureMaxTime         *time.Time `json:"temperatureMaxTime,omitempty"`
	ApparentTemperature        float64    `json:"apparentTemperature,omitempty"`
	ApparentTemperatureMin     float64    `json:"apparentTemperatureMin,omitempty"`
	ApparentTemperatureMinTime *time.Time `json:"apparentTemperatureMinTime,omitempty"`
	ApparentTemperatureMax     float64    `json:"apparentTemperatureMax,omitempty"`
	ApparentTemperatureMaxTime *time.Time `json:"apparentTemperatureMaxTime,omitempty"`
	SnowFall                   float64    `json:"snowFall,omitempty"`
	SnowDepth                  float64    `json:"snowDepth,omitempty"`
	NearestStormBearing        float64    `json:"nearestStormBearing,omitempty"`
	NearestStormDistance       float64    `json:"nearestStormDistance,omitempty"`
	DewPoint                   float64    `json:"dewPoint,omitempty"`
	WindBearing                float64    `json:"windBearing,omitempty"`
	WindChill                  float64    `json:"windchill,omitempty"`
	WindGust                   float64    `json:"windGust,omitempty"`
	WindSpeed                  float64    `json:"windSpeed,omitempty"`
	CloudCover                 float64    `json:"cloudCover,omitempty"`
	Humidity                   float64    `json:"humidity,omitempty"`
	Pressure                   float64    `json:"pressure,omitempty"`
	Visibility                 float64    `json:"visibility,omitempty"`
	Ozone                      float64    `json:"ozone,omitempty"`
	MoonPhase                  float64    `json:"moonPhase,omitempty"`
	MoonRise                   *time.Time `json:"moonRise,omitempty"`
	MoonSet                    *time.Time `json:"moonSet,omitempty"`
	UVIndex                    int64      `json:"uvIndex,omitempty"`
	UVIndexTime                *time.Time `json:"uvIndexTime,omitempty"`
}

// Period .
type Period struct {
	Summary    string        `json:"summary,omitempty"`
	Icon       string        `json:"icon,omitempty"`
	Conditions []*Conditions `json:"data,omitempty"`
}

// Alert .
type Alert struct {
	Title       string   `json:"title,omitempty"`
	Regions     []string `json:"regions,omitempty"`
	Severity    string   `json:"severity,omitempty"`
	Description string   `json:"description,omitempty"`
	Time        int64    `json:"time,omitempty"`
	Expires     float64  `json:"expires,omitempty"`
	URI         string   `json:"uri,omitempty"`
}

// Forecast .
type Forecast struct {
	ID        string      `json:"id,omitempty"`
	Latitude  float64     `json:"latitude,omitempty"`
	Longitude float64     `json:"longitude,omitempty"`
	Timezone  string      `json:"timezone,omitempty"`
	Offset    float64     `json:"offset,omitempty"`
	Current   *Conditions `json:"current,omitempty"`
	Period    *Period     `json:"period,omitempty"`
	Alerts    []*Alert    `json:"alerts,omitempty"`
}

// Feature .
func (f *Forecast) Feature() (*gj.Feature, error) {
	coords := []float64{f.Longitude, f.Latitude}
	t := gj.NewFeature(gj.NewPointGeometry(coords))
	t.ID = f.ID
	t.Properties["current"] = f.Current
	t.Properties["period"] = f.Period
	return t, nil
}

// NewFeatureCollection .
func NewFeatureCollection(forecasts ...*Forecast) (*gj.FeatureCollection, error) {
	fc := gj.NewFeatureCollection()
	if forecasts == nil || len(forecasts) == 0 {
		return fc, nil
	}
	for _, forecast := range forecasts {
		feature, err := forecast.Feature()
		if err != nil {
			return nil, err
		}
		fc.AddFeature(feature)
	}
	return fc, nil
}
