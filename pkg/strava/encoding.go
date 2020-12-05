package strava

import (
	"errors"
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-gpx"
	"github.com/twpayne/go-polyline"

	"github.com/bzimmer/gravl/pkg/common/geo"
)

var _ geo.GPX = &Route{}
var _ geo.GPX = &Streams{}
var _ geo.GPX = &Activity{}

func polylineToLineString(polylines ...string) (*geom.LineString, error) {
	const N = 2
	var coords []float64
	var linestring *geom.LineString
	for _, p := range polylines {
		if p == "" {
			continue
		}
		c, _, err := polyline.DecodeCoords([]byte(p))
		if err != nil {
			return nil, err
		}
		coords = make([]float64, len(c)*N)
		for i := 0; i < len(c); i++ {
			x := N * i
			coords[x+0] = c[i][1]
			coords[x+1] = c[i][0]
		}
		return geom.NewLineStringFlat(geom.XY, coords), nil
	}
	if linestring == nil {
		return nil, errors.New("no valid polyline")
	}
	return linestring, nil
}

func (a *Activity) GPX() (*gpx.GPX, error) {
	ls, err := polylineToLineString(a.Map.Polyline, a.Map.SummaryPolyline)
	if err != nil {
		return nil, err
	}

	mls := geom.NewMultiLineString(ls.Layout())
	err = mls.Push(ls)
	if err != nil {
		return nil, err
	}

	trk := gpx.NewTrkType(mls)
	trk.Src = baseURL

	return &gpx.GPX{
		Metadata: &gpx.MetadataType{
			Name: fmt.Sprintf("%d", a.ID),
			Desc: a.Description,
			Time: a.StartDate,
		},
		Trk: []*gpx.TrkType{trk},
	}, nil
}

func (r *Route) GPX() (*gpx.GPX, error) {
	ls, err := polylineToLineString(r.Map.Polyline, r.Map.SummaryPolyline)
	if err != nil {
		return nil, err
	}

	rte := gpx.NewRteType(ls)
	rte.Src = baseURL

	return &gpx.GPX{
		Metadata: &gpx.MetadataType{
			Name: fmt.Sprintf("%d", r.ID),
			Desc: r.Description,
		},
		Rte: []*gpx.RteType{rte},
	}, nil
}

func (s *Streams) GPX() (*gpx.GPX, error) {
	var dim int
	var layout geom.Layout
	switch {
	case s.Altitude != nil && s.Time != nil:
		dim = 4
		layout = geom.XYZM
	case s.Altitude != nil:
		dim = 3
		layout = geom.XYZ
	case s.Time != nil:
		dim = 3
		layout = geom.XYM
	default:
		dim = 2
		layout = geom.XY
	}

	n := len(s.LatLng.Data)
	coords := make([]float64, dim*n)
	for i := 0; i < n; i++ {
		x := dim * i
		latlng := s.LatLng.Data[i]
		coords[x+0] = latlng[1]
		coords[x+1] = latlng[0]
		switch layout {
		case geom.XYZM:
			coords[x+2] = s.Altitude.Data[i]
			coords[x+3] = s.Time.Data[i]
		case geom.XYZ:
			coords[x+2] = s.Altitude.Data[i]
		case geom.XYM:
			coords[x+2] = s.Time.Data[i]
		case geom.NoLayout, geom.XY:
			// pass
		}
	}

	ls := geom.NewLineStringFlat(layout, coords)
	mls := geom.NewMultiLineString(ls.Layout())
	err := mls.Push(ls)
	if err != nil {
		return nil, err
	}

	trk := gpx.NewTrkType(mls)
	trk.Src = baseURL

	return &gpx.GPX{
		Metadata: &gpx.MetadataType{
			Name: fmt.Sprintf("%d", s.ActivityID),
		},
		Trk: []*gpx.TrkType{trk},
	}, nil
}
