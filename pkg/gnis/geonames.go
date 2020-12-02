package gnis

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/common/geo"
)

const (
	gnisLength = 20
)

// GeoNamesService used to query geonames
type GeoNamesService service

// Query GNIS for geonames
func (s *GeoNamesService) Query(ctx context.Context, state string) ([]*geo.GeographicName, error) {
	uri := fmt.Sprintf("%s/docs/stategaz/%s_Features.zip", baseURL, state)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	ctx = req.Context()
	req.Header.Set("User-Agent", pkg.UserAgent)
	res, err := s.client.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}
	defer res.Body.Close()

	// unfortunately the entire body needs to be read into memory first
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// open the archive
	archive, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}

	// there should be one file
	f := archive.File[0]
	reader, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return parseReader(reader)
}

func parseReader(reader io.Reader) ([]*geo.GeographicName, error) {
	scanner := bufio.NewScanner(reader)

	// skip the header row
	scanner.Scan()

	names := make([]*geo.GeographicName, 0)
	for scanner.Scan() {
		txt := scanner.Text()
		name, err := unmarshal(txt)
		if err != nil {
			log.Error().Err(err).Str("line", txt).Send()
		}
		names = append(names, name)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return names, nil
}

func unmarshal(line string) (*geo.GeographicName, error) { // nolint:gocyclo
	f := &geo.GeographicName{
		Track: geo.Track{
			Source:      baseURL,
			Coordinates: [][]float64{{0, 0, 0}},
		},
	}

	parts := strings.Split(line, "|")
	if len(parts) != gnisLength {
		return nil, fmt.Errorf("found %d parts, expected %d", len(parts), gnisLength)
	}

	for i, s := range parts {
		switch i {
		case 0: // FEATURE_ID
			f.ID = s
		case 1: // FEATURE_NAME
			f.Name = s
		case 2: // FEATURE_CLASS
			f.Class = s
		case 3: // STATE_ALPHA
			f.Locale = s
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
			f.Coordinates[0][1] = x
		case 10: // PRIM_LONG_DEC
			x, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, err
			}
			f.Coordinates[0][0] = x
		case 11: // SOURCE_LAT_DMS
		case 12: // SOURCE_LONG_DMS
		case 13: // SOURCE_LAT_DEC
		case 14: // SOURCE_LONG_DEC
		case 15: // ELEV_IN_M
			if s == "" {
				// not important enough to care about _though_ 0 m elevation is a legit value -- hmmm
				continue
			}
			x, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, err
			}
			f.Coordinates[0][2] = x
		case 16: // ELEV_IN_FT
		case 17: // MAP_NAME
		case 18: // DATE_CREATED
		case 19: // DATE_EDITED
		default:
		}
	}
	return f, nil
}
