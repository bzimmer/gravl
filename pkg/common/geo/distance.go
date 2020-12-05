package geo

import (
	"math"

	"github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid"
	geom "github.com/twpayne/go-geom"
)

// https://github.com/rosshemsley/kalman <= smoothing
// https://github.com/icholy/utm <= lng/lat to x/y

var WGS84 = ellipsoid.Init(
	"WGS84", ellipsoid.Degrees, ellipsoid.Meter,
	ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)

type Step struct {
	Distance2D   float64
	Distance3D   float64
	Bearing      float64
	Difference3D float64
	Grade        float64
	Moving       bool
}

type Summary struct {
	Steps      int
	Distance2D float64
	Distance3D float64
	Ascent     float64
	Descent    float64
}

func To(p, q *geom.Point) Step {
	d2, b := WGS84.To(p.Y(), p.X(), q.Y(), q.X())
	ed := q.Z() - p.Z()
	d3 := math.Sqrt(math.Pow(d2, 2) + math.Pow(ed, 2))
	return Step{
		Distance2D:   d2,
		Distance3D:   d3,
		Bearing:      b,
		Difference3D: ed,
	}
}

func Steps(points ...*geom.Point) []Step {
	if len(points) == 0 {
		return []Step{}
	}
	steps := make([]Step, len(points)-1)
	for i := 0; i < len(points)-1; i++ {
		steps[i] = To(points[i], (points[i+1]))
	}
	return steps
}

func (s *Summary) Add(p Step) {
	s.Steps++
	s.Distance2D += p.Distance2D
	s.Distance3D += p.Distance3D
	if p.Difference3D > 0 {
		s.Ascent += p.Difference3D
	} else {
		s.Descent -= p.Difference3D
	}
}

func Summarize(steps ...Step) Summary {
	s := &Summary{}
	for i := 0; i < len(steps); i++ {
		s.Add(steps[i])
	}
	return *s
}
