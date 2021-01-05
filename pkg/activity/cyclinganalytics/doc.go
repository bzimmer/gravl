// Package cyclinganalytics provides a client to access the Cycling Analytics'
// (API) https://www.cyclinganalytics.com/developer/api.
package cyclinganalytics

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
			"cyclinganalytics.access-token",
			"cyclinganalytics.refresh-token",
			time.Time{}),
		WithAutoRefresh(ctx))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ath, err := client.User.Me(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(ath.Name)
	os.Exit(0)
}
