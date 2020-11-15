package gnis_test

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	gn "github.com/bzimmer/gravl/pkg/gnis"

	"github.com/stretchr/testify/assert"
)

type ZipArchiveTransport struct {
	status   int
	filename string
}

func (t *ZipArchiveTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	filename := "testdata/" + t.filename

	// create the zipfile
	w := &bytes.Buffer{}
	z := zip.NewWriter(w)
	// create the header
	f, err := z.Create(t.filename)
	if err != nil {
		return nil, err
	}
	// copy the contents of the file from disk to the buffer
	datafile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(f, datafile)
	if err != nil {
		return nil, err
	}
	// flush everything
	z.Close()

	header := make(http.Header)
	header.Add("Content-Type", "text/plain")
	header.Add("Content-Encoding", "gzip")

	return &http.Response{
		StatusCode: t.status,
		Body:       ioutil.NopCloser(bytes.NewBuffer(w.Bytes())),
		Header:     header,
	}, nil
}

func Test_Query(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	c, err := gn.NewClient(
		gn.WithTransport(&ZipArchiveTransport{
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
