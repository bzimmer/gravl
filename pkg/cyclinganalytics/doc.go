//go:generate go run ../../cmd/genreadme/genreadme.go

/*
Package cyclinganalytics provides a client to access the Cycling Analytics'
(API) https://www.cyclinganalytics.com/developer/api.

Only the APIs for accessing the logged in user and retrieving Rides has been
implemented.

To use:
	import (
		"context"
		"fmt"
		"time"

		"github.com/bzimmer/gravl/pkg/cyclinganalytics"
	)

	func main() {
		ctx := context.Background()
		client, err := cyclinganalytics.NewClient(
			cyclinganalytics.WithTokenCredentials(
				"cyclinganalytics.access-token",
				"cyclinganalytics.refresh-token",
				time.Time{}),
			cyclinganalytics.WithAutoRefresh(ctx))
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
*/
package cyclinganalytics
