package rwgps

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-gpx"
)

func (t *Trip) GPX() (*gpx.GPX, error) {
	var dim int
	var layout geom.Layout
	switch t.Origin {
	case OriginTrip:
		dim = 4
		layout = geom.XYZM
	case OriginRoute:
		// routes do not have a `time` dimension
		dim = 3
		layout = geom.XYZ
	}

	n := len(t.TrackPoints)
	coords := make([]float64, dim*n)
	for i, tp := range t.TrackPoints {
		x := dim * i
		coords[x+0] = tp.Longitude
		coords[x+1] = tp.Latitude
		switch layout {
		case geom.XYZM:
			coords[x+2] = tp.Elevation
			coords[x+3] = tp.Time
		case geom.XYZ:
			coords[x+2] = tp.Elevation
		case geom.NoLayout, geom.XY, geom.XYM:
			// pass
		}
	}

	x := &gpx.GPX{}
	switch layout {
	case geom.XYZM:
		ls := geom.NewLineStringFlat(layout, coords)
		mls := geom.NewMultiLineString(ls.Layout())
		err := mls.Push(ls)
		if err != nil {
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
