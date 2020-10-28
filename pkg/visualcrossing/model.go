package visualcrossing

import "time"

// Column .
type Column struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
	Unit string `json:"unit"`
}

// Alert .
type Alert struct {
	Event       string `json:"event"`
	Headline    string `json:"headline"`
	Description string `json:"description"`
	Ends        string `json:"ends"`
	Onset       string `json:"onset"`
}

// DateTime .
type DateTime struct {
	time.Time
}

// UnmarshalJSON .
func (t *DateTime) UnmarshalJSON(data []byte) error {
	return nil
}

// Conditions .
type Conditions struct {
	WindDirection    float64  `json:"wdir"`
	Temp             float64  `json:"temp"`
	Sunrise          string   `json:"sunrise"`
	Visibility       float64  `json:"visibility"`
	WindSpeed        float64  `json:"wspd"`
	Icon             string   `json:"icon"`
	Stations         string   `json:"stations"`
	Heatindex        float64  `json:"heatindex"`
	CloudCover       float64  `json:"cloudcover"`
	Precip           float64  `json:"precip"`
	Moonphase        float64  `json:"moonphase"`
	SnowDepth        float64  `json:"snowdepth"`
	SeaLevelPressure float64  `json:"sealevelpressure"`
	Dew              float64  `json:"dew"`
	Sunset           string   `json:"sunset"`
	Humidity         float64  `json:"humidity"`
	WindGust         float64  `json:"wgust"`
	WindChill        float64  `json:"windchill"`
	DateTime         DateTime `json:"datetime"`
}

// Location .
type Location struct {
	StationContributions interface{}  `json:"stationContributions"`
	Conditions           []Conditions `json:"values"`
	ID                   string       `json:"id"`
	Address              string       `json:"address"`
	Name                 string       `json:"name"`
	Index                int          `json:"index"`
	Latitude             float64      `json:"latitude"`
	Longitude            float64      `json:"longitude"`
	Distance             float64      `json:"distance"`
	Time                 float64      `json:"time"`
	Timezone             string       `json:"tz"`
	CurrentConditions    Conditions   `json:"currentConditions"`
	Alerts               []Alert      `json:"alerts"`
}

// Forecast .
type Forecast struct {
	Columns struct {
		WindDirection      Column `json:"wdir"`
		Sunshine           Column `json:"sunshine"`
		Latitude           Column `json:"latitude"`
		CloudCover         Column `json:"cloudcover"`
		Pop                Column `json:"pop"`
		MinTemp            Column `json:"mint"`
		DateTime           Column `json:"datetime"`
		Precip             Column `json:"precip"`
		Dew                Column `json:"dew"`
		Humidity           Column `json:"humidity"`
		Longitude          Column `json:"longitude"`
		Temp               Column `json:"temp"`
		Address            Column `json:"address"`
		MaxTemp            Column `json:"maxt"`
		Visibility         Column `json:"visibility"`
		Wspd               Column `json:"wspd"`
		ResolvedAddress    Column `json:"resolvedAddress"`
		HeatIndex          Column `json:"heatindex"`
		SnowDepth          Column `json:"snowdepth"`
		SeaLevelPressure   Column `json:"sealevelpressure"`
		ShortwaveRadiation Column `json:"sw_radiation"`
		Snow               Column `json:"snow"`
		Name               Column `json:"name"`
		WindGust           Column `json:"wgust"`
		LongwaveRadiation  Column `json:"lw_radiation"`
		Conditions         Column `json:"conditions"`
		WindChill          Column `json:"windchill"`
	} `json:"columns"`
	RemainingCost int         `json:"remainingCost"`
	QueryCost     int         `json:"queryCost"`
	Messages      interface{} `json:"messages"`
	Locations     []Location  `json:"locations"`
}
