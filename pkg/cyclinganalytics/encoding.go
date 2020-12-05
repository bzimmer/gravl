package cyclinganalytics

import (
	"errors"
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-gpx"
)

func (r *Ride) GPX() (*gpx.GPX, error) {
	var dim int
	var layout geom.Layout

	s := r.Streams
	if s.Latitude == nil && s.Longitude == nil {
		return nil, errors.New("required lat and lng not found")
	}
	switch {
	case s.Latitude != nil && s.Longitude != nil && s.Elevation != nil:
		dim = 3
		layout = geom.XYZ
	case s.Latitude != nil && s.Longitude != nil:
		dim = 2
		layout = geom.XY
	}

	n := len(s.Latitude)
	coords := make([]float64, dim*n)
	for i := 0; i < n; i++ {
		x := dim * i
		coords[x+0] = s.Longitude[i]
		coords[x+1] = s.Latitude[i]
		switch layout {
		case geom.XYZ:
			coords[x+2] = s.Elevation[i]
		case geom.NoLayout, geom.XY, geom.XYM, geom.XYZM:
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
			Name: fmt.Sprintf("%d", r.ID),
		},
		Trk: []*gpx.TrkType{trk},
	}, nil
}
