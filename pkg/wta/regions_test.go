package wta_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetRegions(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	client, err := newClient(http.StatusOK, "")
	a.NoError(err)
	regions, err := client.Regions.Regions(context.Background())
	a.NoError(err)
	a.Equal(11, len(regions))
}
