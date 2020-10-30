package noaa

import (
	"time"

	"github.com/bzimmer/gravl/pkg/common/wx"
)

// Fault .
type Fault struct {
	CorrelationID string `json:"correlationId"`
	Title         string `json:"title"`
	Type          string `json:"type"`
	Status        int    `json:"status"`
	Detail        string `json:"detail"`
	Instance      string `json:"instance"`
}

func (f *Fault) Error() string {
	return f.Detail
}

type elevation struct {
	Value    float64 `json:"value"`
	UnitCode string  `json:"unitCode"`
}

type period struct {
	Number           int    `json:"number"`
	Name             string `json:"name"`
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime"`
	IsDaytime        bool   `json:"isDaytime"`
	Temperature      int    `json:"temperature"`
	TemperatureUnit  string `json:"temperatureUnit"`
	TemperatureTrend string `json:"temperatureTrend"`
	WindSpeed        string `json:"windSpeed"`
	WindDirection    string `json:"windDirection"`
	Icon             string `json:"icon"`
	ShortForecast    string `json:"shortForecast"`
	DetailedForecast string `json:"detailedForecast"`
}

type forecast struct {
	Context  []interface{} `json:"@context"`
	Type     string        `json:"type"`
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Updated           time.Time `json:"updated"`
		Units             string    `json:"units"`
		ForecastGenerator string    `json:"forecastGenerator"`
		GeneratedAt       time.Time `json:"generatedAt"`
		UpdateTime        time.Time `json:"updateTime"`
		// formatted as => 2020-10-17T15:00:00+00:00/P7DT10H
		ValidTimes string     `json:"validTimes"`
		Elevation  *elevation `json:"elevation"`
		Periods    []*period  `json:"periods"`
	} `json:"properties"`
}

// WxForecasts .
func (f *forecast) WxForecasts() (*wx.Forecast, error) {
	fcst := &wx.Forecast{
		Longitude: f.Geometry.Coordinates[0][0][0],
		Latitude:  f.Geometry.Coordinates[0][0][1],
		Current:   &wx.Conditions{},
		Period: &wx.Period{
			Conditions: make([]*wx.Conditions, len(f.Properties.Periods)),
		},
	}
	for i, per := range f.Properties.Periods {
		t, err := time.Parse(time.RFC3339, per.StartTime)
		if err != nil {
			return nil, err
		}
		b, err := wx.WindBearing(per.WindDirection)
		if err != nil {
			return nil, err
		}
		fcst.Period.Conditions[i] = &wx.Conditions{
			ValidFrom:   &t,
			Summary:     per.DetailedForecast,
			WindBearing: b,
			Temperature: float64(per.Temperature),
		}
	}
	return fcst, nil
}

// GridPoint .
type GridPoint struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		ID                  string `json:"@id"`
		Type                string `json:"@type"`
		CWA                 string `json:"cwa"`
		ForecastOffice      string `json:"forecastOffice"`
		GridID              string `json:"gridId"`
		GridX               int    `json:"gridX"`
		GridY               int    `json:"gridY"`
		Forecast            string `json:"forecast"`
		ForecastHourly      string `json:"forecastHourly"`
		ForecastGridData    string `json:"forecastGridData"`
		ObservationStations string `json:"observationStations"`
		RelativeLocation    struct {
			Type     string `json:"type"`
			Geometry struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
			Properties struct {
				City     string `json:"city"`
				State    string `json:"state"`
				Distance struct {
					Value    float64 `json:"value"`
					UnitCode string  `json:"unitCode"`
				} `json:"distance"`
				Bearing struct {
					Value    int    `json:"value"`
					UnitCode string `json:"unitCode"`
				} `json:"bearing"`
			} `json:"properties"`
		} `json:"relativeLocation"`
		ForecastZone    string `json:"forecastZone"`
		County          string `json:"county"`
		FireWeatherZone string `json:"fireWeatherZone"`
		TimeZone        string `json:"timeZone"`
		RadarStation    string `json:"radarStation"`
	} `json:"properties"`
}
