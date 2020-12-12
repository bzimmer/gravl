package geo

import geom "github.com/twpayne/go-geom/encoding/geojson"

type GeoJSON interface { // nolint
	GeoJSON() (*geom.Feature, error)
}
