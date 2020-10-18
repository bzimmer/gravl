package gnis

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// GeoNamesService .
type GeoNamesService service

// Query .
func (s *GeoNamesService) Query(ctx context.Context, q string) ([]*Feature, error) {
	return nil, nil
}

func parseFile(filename string) ([]*Feature, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return parseReader(reader)
}

func parseReader(reader io.Reader) ([]*Feature, error) {
	features := make([]*Feature, 0)
	scanner := bufio.NewScanner(reader)

	// skip the header row
	scanner.Scan()

	// now process the data
	for scanner.Scan() {
		txt := scanner.Text()
		feature, err := unmarshal(txt)
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

func unmarshal(line string) (*Feature, error) {
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
				// not important enough to care about _though_ 0 ft elevation is a legit value -- hmmm
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
