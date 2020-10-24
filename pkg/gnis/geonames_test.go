package gnis

import (
	"testing"

	gj "github.com/paulmach/go.geojson"
	"github.com/stretchr/testify/assert"
)

func Test_unmarshall(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	line := "1516141|Barlow Pass|Gap|WA|53|Snohomish|061|480135N|1212638W|48.0264959|-121.4440005|||||721|2365|Bedal|09/10/1979|"

	f, err := unmarshal(line)
	a.Nil(err)
	a.NotNil(f)

	a.Equal(1516141, f.ID)
	a.Equal("Barlow Pass", f.Properties["name"])
	a.Equal("Gap", f.Properties["class"])
	a.Equal("WA", f.Properties["state"])
	a.Equal(-121.4440005, f.Geometry.Point[0])
	a.Equal(48.0264959, f.Geometry.Point[1])
	a.Equal(721.0, f.Geometry.Point[2])
}

func Test_readlines(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	coll, err := parseFile("testdata/WA_Features_20200901.txt")
	a.Nil(err)
	a.NotNil(coll)
	a.Equal(150, len(coll.Features))

	var feature *gj.Feature
	for _, f := range coll.Features {
		if f.Properties["name"] == "The Hump" {
			feature = f
		}
	}
	a.NotNil(feature)
	a.Equal(1527040, feature.ID)
}
