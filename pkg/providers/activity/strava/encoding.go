package strava

import (
	"errors"
	"fmt"
	"time"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-gpx"
	"github.com/twpayne/go-polyline"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/providers/geo"
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

// GPX representation of an activity
func (a *Activity) GPX() (x *gpx.GPX, err error) {
	// minimally require the lat/lng
	if a.Streams != nil && a.Streams.LatLng != nil {
		x, err = toGPXFromStreams(a.Streams, a.StartDate)
		if err != nil {
			return
		}
	} else {
		var ls *geom.LineString
		ls, err = polylineToLineString(a.Map.Polyline, a.Map.SummaryPolyline)
		if err != nil {
			return
		}
		mls := geom.NewMultiLineString(ls.Layout())
		err = mls.Push(ls)
		if err != nil {
			return
		}
		trk := gpx.NewTrkType(mls)
		trk.Src = baseURL
		x = &gpx.GPX{
			Trk: []*gpx.TrkType{trk},
		}
	}
	x.Metadata = &gpx.MetadataType{
		Name: fmt.Sprintf("%d", a.ID),
		Desc: a.Description,
		Time: a.StartDate,
	}
	return
}

// GPX representation of an route
func (r *Route) GPX() (*gpx.GPX, error) {
	ls, err := polylineToLineString(r.Map.Polyline, r.Map.SummaryPolyline)
	if err != nil {
		return nil, err
	}
	rte := gpx.NewRteType(ls)
	rte.Src = baseURL
	return &gpx.GPX{
		Creator: pkg.UserAgent,
		Metadata: &gpx.MetadataType{
			Name: fmt.Sprintf("%d", r.ID),
			Desc: r.Description,
		},
		Rte: []*gpx.RteType{rte},
	}, nil
}

// GPX representation of a streams
func (s *Streams) GPX() (*gpx.GPX, error) {
	return toGPXFromStreams(s, time.Time{})
}

func toGPXFromStreams(s *Streams, start time.Time) (*gpx.GPX, error) {
	if s.LatLng == nil || len(s.LatLng.Data) == 0 {
		return nil, errors.New("missing latlng data")
	}

	var layout geom.Layout
	switch {
	case s.Elevation != nil && s.Time != nil:
		layout = geom.XYZM
	case s.Elevation != nil:
		layout = geom.XYZ
	case s.Time != nil:
		layout = geom.XYM
	default:
		layout = geom.XY
	}

	var toUTC = func(m float64) float64 {
		x := time.Second * time.Duration(m)
		y := start.Add(x)
		return float64(y.Unix())
	}

	n := len(s.LatLng.Data)
	dim := layout.Stride()
	coords := make([]float64, dim*n)
	for i := 0; i < n; i++ {
		x := dim * i
		latlng := s.LatLng.Data[i]
		coords[x+0] = latlng[1]
		coords[x+1] = latlng[0]
		switch layout {
		case geom.XYZM:
			coords[x+2] = s.Elevation.Data[i].Meters()
			coords[x+3] = toUTC(s.Time.Data[i])
		case geom.XYZ:
			coords[x+2] = s.Elevation.Data[i].Meters()
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
		Creator: pkg.UserAgent,
		Metadata: &gpx.MetadataType{
			Name: fmt.Sprintf("%d", s.ActivityID),
		},
		Trk: []*gpx.TrkType{trk},
	}, nil
}
