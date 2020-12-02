package srtm_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/bzimmer/httpwares"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/srtm"
)

type storage struct{}

func (s storage) LoadFile(fn string) ([]byte, error) {
	switch fn {
	case "urls.json", "N48W120.hgt.zip":
		filename := filepath.Join("testdata", fn)
		return ioutil.ReadFile(filename)
	default:
		return []byte{}, fmt.Errorf("missing {%s}", fn)
	}
}

func (s storage) IsNotExists(err error) bool {
	return true
}

func (s storage) SaveFile(fn string, bytes []byte) error {
	return fmt.Errorf("should not be saving files {%s}", fn)
}

func TestClient(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	opts := []srtm.Option{
		srtm.WithStorage(storage{}),
		srtm.WithStorageLocation("testdata")}
	for _, o := range opts {
		ctx := context.Background()
		client, err := srtm.NewClient(
			srtm.WithTransport(&httpwares.TestDataTransport{
				Status:      http.StatusNoContent,
				Filename:    "",
				ContentType: "application/json"}),
			o)
		a.NoError(err)

		elevation, err := client.Elevation.Elevation(ctx, -120, 48)
		a.NoError(err)
		a.Equal(float64(1238), elevation)

		elevations, err := client.Elevation.Elevations(ctx,
			[][]float64{{-120, 48}, {-120, 48.2}})
		a.NoError(err)
		a.Equal([]float64{1238, 1143}, elevations)
	}
}
