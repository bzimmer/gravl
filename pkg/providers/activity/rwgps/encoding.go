package rwgps

import (
	"strconv"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-gpx"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/providers/geo"
)

var _ geo.GPX = (*Trip)(nil)
var _ geo.GeoJSON = (*Trip)(nil)

func (t *Trip) GeoJSON() (*geojson.Feature, error) {
	layout := geom.XYZ
	dim, n := layout.Stride(), len(t.TrackPoints)
	coords := make([]float64, dim*n)
	for i, tp := range t.TrackPoints {
		x := dim * i
		coords[x+0] = tp.Longitude
		coords[x+1] = tp.Latitude
		coords[x+2] = tp.Elevation.Meters()
	}
	// @todo add streams
	g := &geojson.Feature{
		ID:       strconv.FormatInt(t.ID, 10),
		Geometry: geom.NewLineStringFlat(layout, coords),
		Properties: map[string]interface{}{
			"type":   t.Type,
			"name":   t.Name,
			"source": baseURL,
		},
	}
	return g, nil
}

func (t *Trip) GPX() (*gpx.GPX, error) {
	var layout geom.Layout
	switch t.Type {
	case TypeTrip.String():
		layout = geom.XYZM
	case TypeRoute.String():
		// routes do not have a `time` dimension
		layout = geom.XYZ
	}

	n := len(t.TrackPoints)
	dim := layout.Stride()
	coords := make([]float64, dim*n)
	for i, tp := range t.TrackPoints {
		x := dim * i
		coords[x+0] = tp.Longitude
		coords[x+1] = tp.Latitude
		switch layout {
		case geom.XYZM:
			coords[x+3] = tp.Time
			fallthrough
		case geom.XYZ:
			coords[x+2] = tp.Elevation.Meters()
		case geom.NoLayout, geom.XY, geom.XYM:
			// pass
		}
	}

	x := &gpx.GPX{
		Creator: pkg.UserAgent,
		Metadata: &gpx.MetadataType{
			Name: strconv.FormatInt(t.ID, 10),
		},
	}
	switch layout {
	case geom.XYZM:
		ls := geom.NewLineStringFlat(layout, coords)
		mls := geom.NewMultiLineString(ls.Layout())
		if err := mls.Push(ls); err != nil {
			return nil, err
		}
		x.Trk = []*gpx.TrkType{gpx.NewTrkType(mls)}
	case geom.XYZ:
		ls := geom.NewLineStringFlat(layout, coords)
		x.Rte = []*gpx.RteType{gpx.NewRteType(ls)}
	case geom.NoLayout, geom.XY, geom.XYM:
		// pass
	}
	return x, nil
}
