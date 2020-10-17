package gnis

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_unmarshall(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	line := "1516141|Barlow Pass|Gap|WA|53|Snohomish|061|480135N|1212638W|48.0264959|-121.4440005|||||721|2365|Bedal|09/10/1979|"

	g := New()
	f, err := g.unmarshal(line)
	a.Nil(err)
	a.NotNil(f)

	a.Equal(1516141, f.ID)
	a.Equal("Barlow Pass", f.Name)
	a.Equal("Gap", f.Class)
	a.Equal("WA", f.State)
	a.Equal(-121.4440005, f.Longitude)
	a.Equal(48.0264959, f.Latitude)
	a.Equal(2365, f.Elevation)
}

func Test_readlines(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	dir, _ := os.Getwd()
	filename := filepath.Join(dir, "../../testdata", "WA_Features_20200901.txt")

	g := New()
	features, err := g.ParseFile(filename)
	a.Nil(err)
	a.NotNil(features)

	a.Equal(150, len(features))

	var feature *Feature
	for _, f := range features {
		if f.Name == "The Hump" {
			feature = f
		}
	}
	a.NotNil(feature)
	a.Equal(1527040, feature.ID)
}
