package gravl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/commands/internal"
)

/*
activities_file=$(mktemp /tmp/gravl.XXXXXXXXX)
buntdb_file=$(mktemp /tmp/gravl.XXXXXXXXX)

gravl -c strava athlete
gravl -c strava activities -N 25 > $activities_file
gravl store update --input "file://$activities_file" --output "bunt://$buntdb_file"
gravl pass --input "bunt://$buntdb_file"
gravl -c store export --input "bunt://$buntdb_file"

act_id=$(head -1 $activities_file | jq ".id")
gravl -c strava activity $act_id

gravl -c rwgps athlete
gravl -c rwgps activities > $activities_file
act_id=$(head -1 $activities_file | jq ".id")
gravl -c rwgps activity $act_id

gravl -c cyclinganalytics athlete
gravl -c cyclinganalytics activities > $activities_file
act_id=$(head -1 $activities_file | jq ".id")
gravl -c cyclinganalytics activity $act_id
*/

func TestGravlIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	a := assert.New(t)

	c := internal.Gravl()
	s := <-c.Start()
	a.Equal(0, s.Exit)

	c = internal.Gravl("foo", "bar", "baz")
	s = <-c.Start()
	a.Equal(1, s.Exit)
}
