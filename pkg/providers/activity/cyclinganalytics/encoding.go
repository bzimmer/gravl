package cyclinganalytics

import (
	"errors"
	"strconv"
	"time"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-gpx"

	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/providers/geo"
)

var _ geo.GPX = (*Ride)(nil)

func (r *Ride) GPX() (*gpx.GPX, error) {
	var layout geom.Layout

	s := r.Streams
	if s.Latitude == nil && s.Longitude == nil {
		return nil, errors.New("required lat and lng not found")
	}
	switch {
	case s.Latitude != nil && s.Longitude != nil && s.Elevation != nil:
		layout = geom.XYZM
	case s.Latitude != nil && s.Longitude != nil:
		layout = geom.XYM
	}

	n := len(s.Latitude)
	dim := layout.Stride()
	coords := make([]float64, dim*n)
	for i := 0; i < n; i++ {
		x := dim * i
		t := float64(r.LocalDatetime.Time.Add(time.Second * time.Duration(i)).Unix())

		coords[x+0] = s.Longitude[i]
		coords[x+1] = s.Latitude[i]
		switch layout {
		case geom.XYZM:
			coords[x+2] = s.Elevation[i]
			coords[x+3] = t
		case geom.XYM:
			coords[x+2] = t
		case geom.NoLayout, geom.XY, geom.XYZ:
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
			Name: strconv.FormatInt(r.ID, 10),
		},
		Trk: []*gpx.TrkType{trk},
	}, nil
}
