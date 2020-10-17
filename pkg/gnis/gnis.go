package gnis

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	gnisLength = 20
	baseURL    = "https://geonames.usgs.gov/docs/stategaz/%s_Features.zip"
)

// Feature .
type Feature struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Class     string  `json:"class"`
	State     string  `json:"state"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation int     `json:"elevation"`
}

// GNIS .
type GNIS struct{}

// New
func New() *GNIS {
	return &GNIS{}
}

func (g *GNIS) ParseFile(filename string) ([]*Feature, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return g.ParseReader(reader)
}

func (g *GNIS) ParseReader(reader io.Reader) ([]*Feature, error) {
	features := make([]*Feature, 0)
	scanner := bufio.NewScanner(reader)

	// skip the header row
	scanner.Scan()

	// now process the data
	for scanner.Scan() {
		txt := scanner.Text()
		feature, err := g.unmarshal(txt)
		if err != nil {
			log.Error().Err(err).Str("line", txt).Send()
		}
		features = append(features, feature)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return features, nil
}

func (g *GNIS) unmarshal(line string) (*Feature, error) {
	f := &Feature{}

	parts := strings.Split(line, "|")
	if len(parts) != gnisLength {
		return nil, fmt.Errorf("found %d parts, expected %d", len(parts), gnisLength)
	}

	for i, s := range parts {
		switch i {
		case 0: // FEATURE_ID
			x, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			f.ID = x
		case 1: // FEATURE_NAME
			f.Name = s
		case 2: // FEATURE_CLASS
			f.Class = s
		case 3: // STATE_ALPHA
			f.State = s
		case 4: // STATE_NUMERIC
		case 5: // COUNTY_NAME
		case 6: // COUNTY_NUMERIC
		case 7: // PRIMARY_LAT_DMS
		case 8: // PRIM_LONG_DMS
		case 9: // PRIM_LAT_DEC
			x, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, err
			}
			f.Latitude = x
		case 10: // PRIM_LONG_DEC
			x, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, err
			}
			f.Longitude = x
		case 11: // SOURCE_LAT_DMS
		case 12: // SOURCE_LONG_DMS
		case 13: // SOURCE_LAT_DEC
		case 14: // SOURCE_LONG_DEC
		case 15: // ELEV_IN_M
		case 16: // ELEV_IN_FT
			if s == "" {
				continue
			}
			x, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			f.Elevation = x
		case 17: // MAP_NAME
		case 18: // DATE_CREATED
		case 19: // DATE_EDITED
		default:
		}
	}
	return f, nil
}
