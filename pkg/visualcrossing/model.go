package visualcrossing

import "time"

// https://www.visualcrossing.com/resources/documentation/weather-api/weather-api-json-result-structure/

// Fault .
type Fault struct {
	ErrorCode     int    `json:"errorCode"`
	ExecutionTime int    `json:"executionTime"`
	Message       string `json:"message"`
	SessionID     string `json:"sessionId"`
}

func (f Fault) Error() string {
	return f.Message
}

// Column .
type Column struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
	Unit string `json:"unit"`
}

// Columns .
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

// Alert .
type Alert struct {
	Event       string    `json:"event"`
	Headline    string    `json:"headline"`
	Description string    `json:"description"`
	Ends        time.Time `json:"ends"`
	Onset       time.Time `json:"onset"`
}

// Conditions .
type Conditions struct {
	CloudCover       float64   `json:"cloudcover"`
	DateTime         time.Time `json:"datetimeStr"`
	Dew              float64   `json:"dew"`
	HeatIndex        float64   `json:"heatindex"`
	Humidity         float64   `json:"humidity"`
	MoonPhase        float64   `json:"moonphase"`
	Precip           float64   `json:"precip"`
	SeaLevelPressure float64   `json:"sealevelpressure"`
	SnowDepth        float64   `json:"snowdepth"`
	Stations         string    `json:"stations"`
	Sunrise          time.Time `json:"sunrise"`
	Sunset           time.Time `json:"sunset"`
	Temp             float64   `json:"temp"`
	Visibility       float64   `json:"visibility"`
	WindChill        float64   `json:"windchill"`
	WindDirection    float64   `json:"wdir"`
	WindGust         float64   `json:"wgust"`
	WindSpeed        float64   `json:"wspd"`
}

// Location .
type Location struct {
	Address           string       `json:"address"`
	Alerts            []Alert      `json:"alerts"`
	Conditions        []Conditions `json:"values"`
	CurrentConditions Conditions   `json:"currentConditions"`
	Distance          float64      `json:"distance"`
	ID                string       `json:"id"`
	Index             int          `json:"index"`
	Latitude          float64      `json:"latitude"`
	Longitude         float64      `json:"longitude"`
	Name              string       `json:"name"`
	Time              float64      `json:"time"`
	Timezone          string       `json:"tz"`
	// StationContributions interface{}  `json:"stationContributions"`
}

// Forecast .
type Forecast struct {
	Columns       Columns     `json:"columns"`
	RemainingCost int         `json:"remainingCost"`
	QueryCost     int         `json:"queryCost"`
	Messages      interface{} `json:"messages"`
	Locations     []Location  `json:"locations"`
}
