package gnis

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	gj "github.com/paulmach/go.geojson"
	"github.com/stretchr/testify/assert"
)

type ZipArchiveTransport struct {
	status   int
	filename string
}

func (ar *ZipArchiveTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	dir, _ := os.Getwd()
	filename := filepath.Join(dir, "../../testdata", ar.filename)

	// create the zipfile
	w := &bytes.Buffer{}
	z := zip.NewWriter(w)
	// create the header
	f, err := z.Create(ar.filename)
	if err != nil {
		return nil, err
	}
	// copy the contents of the file from disk to the buffer
	datafile, err := os.Open(filename)
	_, err = io.Copy(f, datafile)
	if err != nil {
		return nil, err
	}
	// flush everything
	z.Close()

	return &http.Response{
		StatusCode: ar.status,
		Body:       ioutil.NopCloser(bytes.NewBuffer(w.Bytes())),
		Header:     make(http.Header),
	}, nil
}

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

	dir, _ := os.Getwd()
	filename := filepath.Join(dir, "../../testdata", "WA_Features_20200901.txt")

	coll, err := parseFile(filename)
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

func Test_Query(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := NewClient(
		WithTransport(&ZipArchiveTransport{
			status:   http.StatusOK,
			filename: "WA_Features_20200901.txt",
		}),
	)
	a.NoError(err)
	a.NotNil(c)

	b := context.Background()
	coll, err := c.GeoNames.Query(b, "WA")
	a.NoError(err)
	a.NotNil(coll)
	a.Equal(150, len(coll.Features))
	a.Equal("Blue Buck Ridge", coll.Features[109].Properties["name"])
}
