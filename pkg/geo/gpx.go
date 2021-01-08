package geo

import (
	"math"
	"time"

	"github.com/golang/geo/s2"
	geom "github.com/twpayne/go-geom"
	gpx "github.com/twpayne/go-gpx"
)

const earthRadiusM = 6367000.0
const movingTimeThreshold = 0.1

type GPX interface {
	GPX() (*gpx.GPX, error)
}

type Summary struct {
	Tracks      int           `json:"tracks,omitempty"`
	Routes      int           `json:"routes,omitempty"`
	Segments    int           `json:"segments,omitempty"`
	Points      int           `json:"points,omitempty"`
	Distance2D  float64       `json:"distance2d,omitempty"`
	Distance3D  float64       `json:"distance3d,omitempty"`
	Ascent      float64       `json:"ascent,omitempty"`
	Descent     float64       `json:"descent,omitempty"`
	StartTime   time.Time     `json:"start_time,omitempty"`
	MovingTime  time.Duration `json:"moving_time,omitempty"`
	StoppedTime time.Duration `json:"stopped_time,omitempty"`
}

func (s Summary) TotalTime() time.Duration {
	return s.MovingTime + s.StoppedTime
}

func distance(p, q *gpx.WptType) float64 {
	llp := s2.LatLngFromDegrees(p.Lat, p.Lon)
	llq := s2.LatLngFromDegrees(q.Lat, q.Lon)
	return float64(llp.Distance(llq)) * earthRadiusM
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

func SummarizeTracks(gpx *gpx.GPX) Summary {
	s := Summary{
		MovingTime:  0 * time.Second,
		StoppedTime: 0 * time.Second,
	}
	for _, track := range gpx.Trk {
		s.Tracks++
		for _, segment := range track.TrkSeg {
			n := len(segment.TrkPt)
			s.Segments++
			for j, point := range segment.TrkPt {
				s.Points++
				if j == 0 {
					if (s.StartTime == time.Time{}) || point.Time.Before(s.StartTime) {
						s.StartTime = point.Time
					}
				}
				if j < n-1 {
					p, q := point, segment.TrkPt[j+1]
					d2 := distance(p, q)
					ele := q.Ele - p.Ele
					d3 := math.Sqrt(math.Pow(d2, 2) + math.Pow(ele, 2))

					s.Distance2D += d2
					s.Distance3D += d3
					switch {
					case ele > 0:
						s.Ascent += ele
					case ele < 0:
						s.Descent -= ele
					}

					t := q.Time.Sub(p.Time)
					if d2 > movingTimeThreshold {
						s.MovingTime += t
					} else {
						s.StoppedTime += t
					}
				}
			}
		}
	}
	return s
}

func SummarizeRoutes(gpx *gpx.GPX) Summary {
	s := Summary{
		MovingTime:  0 * time.Second,
		StoppedTime: 0 * time.Second,
	}
	for _, rte := range gpx.Rte {
		s.Routes++
		n := len(rte.RtePt)
		for j, point := range rte.RtePt {
			s.Points++
			if j < n-1 {
				p, q := point, rte.RtePt[j+1]
				d2 := distance(p, q)
				elv := q.Ele - p.Ele
				d3 := math.Sqrt(math.Pow(d2, 2) + math.Pow(elv, 2))

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
	}
	return s
}
