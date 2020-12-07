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
	Distance2D   float64 `json:"distance2d"`
	Distance3D   float64 `json:"distance3d"`
	Bearing      float64 `json:"bearing"`
	Difference3D float64 `json:"difference3d"`
	Grade        float64 `json:"grade"`
	Moving       bool    `json:"moving"`
}

type Summary struct {
	Steps      int     `json:"steps"`
	Distance2D float64 `json:"distance2d"`
	Distance3D float64 `json:"distance3d"`
	Ascent     float64 `json:"ascent"`
	Descent    float64 `json:"descent"`
}

// To returns the difference between two points
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

// Steps returns a slice steps from one point to another
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

// Add adds a Step to the summary
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

// Summarize the steps into a Summary
func Summarize(steps ...Step) Summary {
	s := &Summary{}
	for i := 0; i < len(steps); i++ {
		s.Add(steps[i])
	}
	return *s
}
