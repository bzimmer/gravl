package cyclinganalytics_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRides(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusOK, "me-rides.json")
	a.NoError(err)
	ctx := context.Background()
	rides, err := client.Rides.Rides(ctx)
	a.NoError(err)
	a.NotNil(rides)
	a.Equal(2, len(rides))
}
