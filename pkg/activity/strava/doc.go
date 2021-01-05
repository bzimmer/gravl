// Package strava provides a client to access the Strava API.

// Information about the Strava API can be found at Strava's (site)
// https://developers.com/. This client provides read-only
// access and is not intended to be a complete implementation.
package strava

import (
	"context"
	"fmt"
	"os"
	"time"
)

func Example() {
	ctx := context.Background()
	client, err := NewClient(
		WithTokenCredentials(
			"access-token",
			"refresh-token",
			// set the expiry in the past to force token refreshing
			time.Now().Add(-1*time.Minute)),
		WithClientCredentials(
			"client-id",
			"client-secret"),
		// with auto refresh enabled an expired access token will be refreshed using the refresh token provided
		WithAutoRefresh(ctx))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	act, err := client.Activity.Activity(ctx, 4273918760)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(act)
	os.Exit(0)
}
