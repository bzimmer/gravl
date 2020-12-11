package geo

import (
	"math"
	"time"

	"github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid"
	geom "github.com/twpayne/go-geom"
	gpx "github.com/twpayne/go-gpx"
)

const movingTimeThreshold = 0.1

type GPX interface {
	GPX() (*gpx.GPX, error)
}

type Summary struct {
	Tracks      int           `json:"tracks"`
	Routes      int           `json:"routes"`
	Segments    int           `json:"segments"`
	Points      int           `json:"points"`
	Distance2D  float64       `json:"distance2d"`
	Distance3D  float64       `json:"distance3d"`
	Ascent      float64       `json:"ascent"`
	Descent     float64       `json:"descent"`
	StartTime   time.Time     `json:"start_time"`
	MovingTime  time.Duration `json:"moving_time"`
	StoppedTime time.Duration `json:"stopped_time"`
}

func (s Summary) TotalTime() time.Duration {
	return s.MovingTime + s.StoppedTime
}

var WGS84 = ellipsoid.Init(
	"WGS84", ellipsoid.Degrees, ellipsoid.Meter,
	ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)

func Flatten(gpx *gpx.GPX, layout geom.Layout) *geom.LineString {
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

func Summarize(gpx *gpx.GPX) Summary {
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
					d2, _ := WGS84.To(p.Lat, p.Lon, q.Lat, q.Lon)
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
