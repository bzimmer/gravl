package geo

import "github.com/twpayne/go-geom/encoding/geojson"

type GeoJSON interface { // nolint
	GeoJSON() (*geojson.Feature, error)
}
