package strava_test

import (
	"context"
	"math"
	"net/http"
	"testing"

	"github.com/martinlindhe/unit"
	"github.com/stretchr/testify/assert"
)

func Test_Athlete(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	client, err := newClient(http.StatusOK, "athlete.json")
	a.NoError(err)
	ctx := context.Background()
	ath, err := client.Athlete.Athlete(ctx)
	a.NoError(err, "failed decoding")
	a.NotNil(ath)
	a.Equal(1122, ath.ID)
	a.Equal(1, len(ath.Bikes))
	a.Equal(1, len(ath.Shoes))
}

func Test_AthleteNotAuthorized(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	client, err := newClient(http.StatusUnauthorized, "athlete_unauthorized.json")
	a.NoError(err)
	ctx := context.Background()
	ath, err := client.Athlete.Athlete(ctx)
	a.Error(err)
	a.Nil(ath)
}

func Test_AthleteStats(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	client, err := newClient(http.StatusOK, "athlete_stats.json")
	a.NoError(err)
	ctx := context.Background()
	sts, err := client.Athlete.Stats(ctx, 88273)
	a.NoError(err, "failed decoding")
	a.NotNil(sts)
	a.Equal(float64(14492298), math.Trunc(sts.AllRideTotals.Distance.Meters()))
	a.Equal(unit.Duration(12441), sts.AllSwimTotals.ElapsedTime)
	a.Equal(float64(1597), math.Trunc(sts.BiggestClimbElevationGain.Meters()))
}
