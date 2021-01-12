package geo_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	gpx "github.com/twpayne/go-gpx"

	"github.com/bzimmer/gravl/pkg/providers/geo"
)

func TestSummarize(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	x, err := gpx.Read(bytes.NewBufferString(utrecht))

	a.NoError(err)
	a.NotNil(x)

	s := geo.SummarizeTracks(x)
	a.Equal(13, s.Points, "points")
	a.Equal(10.0, s.Ascent, "ascent")
	a.Equal(15.0, s.Descent, "descent")
	a.Equal(14683, int(s.Distance2D), "distance")
	a.Equal(14683, int(s.Distance2D), "distance")
}

var utrecht = `
<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1"
	creator="GPS Visualizer Sandbox https://www.gpsvisualizer.com/draw/"
	xmlns="http://www.topografix.com/GPX/1/1"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd">
	<trk>
		<name></name>
		<desc>Length: 14.69 km (9.129 mi)</desc>
		<trkseg>
			<trkpt lat="52.1051345" lon="5.1198006">
				<ele>5</ele>
			</trkpt>
			<trkpt lat="52.105978" lon="5.136795"></trkpt>
			<trkpt lat="52.0957497" lon="5.1486397"></trkpt>
			<trkpt lat="52.0866793" lon="5.1539612">
				<ele>10</ele>
			</trkpt>
			<trkpt lat="52.078451" lon="5.1568794">
				<ele>9</ele>
			</trkpt>
			<trkpt lat="52.0682163" lon="5.137825"></trkpt>
			<trkpt lat="52.0653671" lon="5.1163673"></trkpt>
			<trkpt lat="52.0705378" lon="5.0992012"></trkpt>
			<trkpt lat="52.079295" lon="5.0895882"></trkpt>
			<trkpt lat="52.0903709" lon="5.08564"></trkpt>
			<trkpt lat="52.0998624" lon="5.0861549"></trkpt>
			<trkpt lat="52.1039747" lon="5.0979996"></trkpt>
			<trkpt lat="52.1062943" lon="5.1146507"></trkpt>
		</trkseg>
	</trk>
</gpx>
`
