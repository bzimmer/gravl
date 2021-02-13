package geo

import (
	"context"
	"math"
	"time"

	"github.com/golang/geo/s2"
	"github.com/martinlindhe/unit"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-gpx"
)

const (
	earthRadiusM        = 6367000.0
	movingTimeThreshold = 0.1
)

// GPX instances can return a gpx instance
type GPX interface {
	// GPX returns a gpx instance
	GPX() (*gpx.GPX, error)
}

// Elevator instances can provide an elevation for a point
type Elevator interface {
	// Elevation for a point
	Elevation(ctx context.Context, point *geom.Point) (float64, error)
}

type Summary struct {
	Filename    string        `json:"filename,omitempty"`
	Tracks      int           `json:"tracks,omitempty"`
	Routes      int           `json:"routes,omitempty"`
	Segments    int           `json:"segments,omitempty"`
	Points      int           `json:"points,omitempty"`
	Waypoints   int           `json:"waypoints,omitempty"`
	Distance2D  unit.Length   `json:"distance2d,omitempty"`
	Distance3D  unit.Length   `json:"distance3d,omitempty"`
	Ascent      unit.Length   `json:"ascent,omitempty"`
	Descent     unit.Length   `json:"descent,omitempty"`
	StartTime   time.Time     `json:"start_time,omitempty"`
	MovingTime  unit.Duration `json:"moving_time,omitempty"`
	StoppedTime unit.Duration `json:"stopped_time,omitempty"`
}

func (s Summary) TotalTime() unit.Duration {
	return s.MovingTime + s.StoppedTime
}

func distance(p, q *gpx.WptType) unit.Length {
	llp := s2.LatLngFromDegrees(p.Lat, p.Lon)
	llq := s2.LatLngFromDegrees(q.Lat, q.Lon)
	return unit.Length(llp.Distance(llq)) * earthRadiusM
}

func CorrectElevations(ctx context.Context, gpx *gpx.GPX, elevator Elevator) error {
	for _, track := range gpx.Trk {
		for _, segment := range track.TrkSeg {
			for _, point := range segment.TrkPt {
				elv, err := elevator.Elevation(ctx, point.Geom(geom.XY))
				if err != nil {
					return err
				}
				// fmt.Println(point.Lat, point.Lon, point.Ele, elv)
				// log.Debug().
				// 	Float64("lat", point.Lat).
				// 	Float64("lon", point.Lon).
				// 	Float64("elv", point.Ele).
				// 	Float64("srtm", elv).
				// 	Msg("track")
				point.Ele = elv
			}
		}
	}
	for _, rte := range gpx.Rte {
		for _, point := range rte.RtePt {
			elv, err := elevator.Elevation(ctx, point.Geom(geom.XY))
			if err != nil {
				return err
			}
			point.Ele = elv
		}
	}
	return nil
}

func FlattenTracks(gpx *gpx.GPX, layout geom.Layout) *geom.LineString {
	var coords []float64
	for _, track := range gpx.Trk {
		for _, segment := range track.TrkSeg {
			for _, point := range segment.TrkPt {
				c := point.Geom(layout).FlatCoords()
				for i := 0; i < len(c); i++ {
					coords = append(coords, c[i])
				}
			}
		}
	}
	return geom.NewLineStringFlat(layout, coords)
}

func FlattenRoutes(gpx *gpx.GPX, layout geom.Layout) *geom.LineString {
	var coords []float64
	for _, rte := range gpx.Rte {
		for _, point := range rte.RtePt {
			c := point.Geom(layout).FlatCoords()
			for i := 0; i < len(c); i++ {
				coords = append(coords, c[i])
			}
		}
	}
	return geom.NewLineStringFlat(layout, coords)
}

func SummarizeTracks(gpx *gpx.GPX) *Summary {
	s := &Summary{Waypoints: len(gpx.Wpt)}
	for _, track := range gpx.Trk {
		s.Tracks++
		for _, segment := range track.TrkSeg {
			n := len(segment.TrkPt)
			s.Segments++
			for j, point := range segment.TrkPt {
				s.Points++
				if j == n-1 {
					continue
				}
				if j == 0 {
					if (s.StartTime == time.Time{}) || point.Time.Before(s.StartTime) {
						s.StartTime = point.Time
					}
				}
				p, q := point, segment.TrkPt[j+1]
				d2 := distance(p, q)
				ele := unit.Length(q.Ele - p.Ele)
				d3 := unit.Length(math.Sqrt(math.Pow(d2.Meters(), 2) + math.Pow(ele.Meters(), 2)))

				s.Distance2D += d2
				s.Distance3D += d3
				switch {
				case ele > 0:
					s.Ascent += ele
				case ele < 0:
					s.Descent -= ele
				}

				t := unit.Duration(q.Time.Sub(p.Time).Seconds())
				if d2 > movingTimeThreshold {
					s.MovingTime += t
				} else {
					s.StoppedTime += t
				}
			}
		}
	}
	return s
}

func SummarizeRoutes(gpx *gpx.GPX) *Summary {
	s := &Summary{Waypoints: len(gpx.Wpt)}
	for _, rte := range gpx.Rte {
		s.Routes++
		n := len(rte.RtePt)
		for j, point := range rte.RtePt {
			s.Points++
			if j == n-1 {
				continue
			}
			p, q := point, rte.RtePt[j+1]
			d2 := distance(p, q)
			elv := unit.Length(q.Ele - p.Ele)
			d3 := unit.Length(math.Sqrt(math.Pow(d2.Meters(), 2) + math.Pow(elv.Meters(), 2)))

			s.Distance2D += d2
			s.Distance3D += d3
			switch {
			case elv > 0:
				s.Ascent += elv
			case elv < 0:
				s.Descent -= elv
			}
		}
	}
	return s
}
