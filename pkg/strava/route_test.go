package strava_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bzimmer/gravl/pkg/common/route"
	"github.com/stretchr/testify/assert"
)

func Test_Route(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusOK, "route.json")
	ctx := context.Background()
	rte, err := client.Route.Route(ctx, 26587226)
	a.NoError(err)
	a.NotNil(rte)
	a.Equal("26587226", rte.ID)
	a.Equal(2076, len(rte.Coordinates))
	a.Equal(route.Planned, rte.Origin)
}
