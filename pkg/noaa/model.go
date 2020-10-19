package noaa

import (
	"encoding/xml"
	"time"
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

// Elevation .
type Elevation struct {
	Value    float64 `json:"value"`
	UnitCode string  `json:"unitCode"`
}

// Period .
type Period struct {
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

// Forecast .
type Forecast struct {
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
		Elevation  *Elevation `json:"elevation"`
		Periods    []*Period  `json:"periods"`
	} `json:"properties"`
}

// GridPoint .
type GridPoint struct {
	Context  []interface{} `json:"@context"`
	ID       string        `json:"id"`
	Type     string        `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		ID                  string `json:"@id"`
		Type                string `json:"@type"`
		Cwa                 string `json:"cwa"`
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

// DWML .
// https://vlab.ncep.noaa.gov/web/mdl/ndfd-web-services
// https://schemas.liquid-technologies.com/DWML/0/
type DWML struct {
	XMLName                   xml.Name `xml:"dwml"`
	Text                      string   `xml:",chardata"`
	Version                   string   `xml:"version,attr"`
	XSD                       string   `xml:"xsd,attr"`
	XSI                       string   `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	Head                      struct {
		Text    string `xml:",chardata"`
		Product struct {
			Text            string `xml:",chardata"`
			ConciseName     string `xml:"concise-name,attr"`
			OperationalMode string `xml:"operational-mode,attr"`
			SrsName         string `xml:"srsName,attr"`
			CreationDate    struct {
				Text             string `xml:",chardata"`
				RefreshFrequency string `xml:"refresh-frequency,attr"`
			} `xml:"creation-date"`
		} `xml:"product"`
		Source struct {
			Text             string `xml:",chardata"`
			ProductionCenter string `xml:"production-center"`
			Credit           string `xml:"credit"`
			MoreInformation  string `xml:"more-information"`
		} `xml:"source"`
	} `xml:"head"`
	Data struct {
		Text     string `xml:",chardata"`
		Location struct {
			Text        string `xml:",chardata"`
			LocationKey string `xml:"location-key"`
			Point       struct {
				Text      string `xml:",chardata"`
				Latitude  string `xml:"latitude,attr"`
				Longitude string `xml:"longitude,attr"`
			} `xml:"point"`
			AreaDescription string `xml:"area-description"`
			Height          struct {
				Text        string `xml:",chardata"`
				Datum       string `xml:"datum,attr"`
				HeightUnits string `xml:"height-units,attr"`
			} `xml:"height"`
		} `xml:"location"`
		MoreWeatherInformation struct {
			Text               string `xml:",chardata"`
			ApplicableLocation string `xml:"applicable-location,attr"`
		} `xml:"moreWeatherInformation"`
		TimeLayout struct {
			Text           string   `xml:",chardata"`
			TimeCoordinate string   `xml:"time-coordinate,attr"`
			Summarization  string   `xml:"summarization,attr"`
			LayoutKey      string   `xml:"layout-key"`
			StartValidTime []string `xml:"start-valid-time"`
			EndValidTime   []string `xml:"end-valid-time"`
		} `xml:"time-layout"`
		Parameters struct {
			Text               string `xml:",chardata"`
			ApplicableLocation string `xml:"applicable-location,attr"`
			Temperature        []struct {
				Text       string `xml:",chardata"`
				Type       string `xml:"type,attr"`
				TimeLayout string `xml:"time-layout,attr"`
				Value      []struct {
					Text string `xml:",chardata"`
					Nil  string `xml:"nil,attr"`
				} `xml:"value"`
			} `xml:"temperature"`
			WindSpeed []struct {
				Text       string `xml:",chardata"`
				Type       string `xml:"type,attr"`
				TimeLayout string `xml:"time-layout,attr"`
				Value      []struct {
					Text string `xml:",chardata"`
					Nil  string `xml:"nil,attr"`
				} `xml:"value"`
			} `xml:"wind-speed"`
			CloudAmount struct {
				Text       string   `xml:",chardata"`
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Value      []string `xml:"value"`
			} `xml:"cloud-amount"`
			ProbabilityOfPrecipitation struct {
				Text       string   `xml:",chardata"`
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Value      []string `xml:"value"`
			} `xml:"probability-of-precipitation"`
			Humidity struct {
				Text       string `xml:",chardata"`
				Type       string `xml:"type,attr"`
				Units      string `xml:"units,attr"`
				TimeLayout string `xml:"time-layout,attr"`
				Value      []struct {
					Text string `xml:",chardata"`
					Nil  string `xml:"nil,attr"`
				} `xml:"value"`
			} `xml:"humidity"`
			Direction struct {
				Text       string   `xml:",chardata"`
				Type       string   `xml:"type,attr"`
				Units      string   `xml:"units,attr"`
				TimeLayout string   `xml:"time-layout,attr"`
				Value      []string `xml:"value"`
			} `xml:"direction"`
			HourlyQPF struct {
				Text       string `xml:",chardata"`
				Type       string `xml:"type,attr"`
				Units      string `xml:"units,attr"`
				TimeLayout string `xml:"time-layout,attr"`
				Value      []struct {
					Text string `xml:",chardata"`
					Nil  string `xml:"nil,attr"`
				} `xml:"value"`
			} `xml:"hourly-qpf"`
			Weather struct {
				Text              string `xml:",chardata"`
				TimeLayout        string `xml:"time-layout,attr"`
				WeatherConditions []struct {
					Text  string `xml:",chardata"`
					Nil   string `xml:"nil,attr"`
					Value []struct {
						Text        string `xml:",chardata"`
						WeatherType string `xml:"weather-type,attr"`
						Coverage    string `xml:"coverage,attr"`
						Additive    string `xml:"additive,attr"`
					} `xml:"value"`
				} `xml:"weather-conditions"`
			} `xml:"weather"`
		} `xml:"parameters"`
	} `xml:"data"`
}
