# strava

Package strava provides a client to access the Strava API.

Information about the Strava API can be found at Strava's (site)
[https://developers.strava.com/](https://developers.strava.com/). This client provides read-only
access and is not intended to be a complete implementation.

To use:

```go
import (
	"context"
	"fmt"
	"os"

	"github.com/bzimmer/gravl/pkg/strava"
)

func main() {
	client, err := strava.NewClient(
		strava.WithTokenCredentials(
			"strava.access-token",
			"strava.refresh-token",
			// set the expiry in the past to force token refreshing
			time.Now().Add(-1*time.Minute)),
		strava.WithClientCredentials(
			"strava.client-id",
			"strava.client-secret"),
		// with auto refresh enabled an expired access token will be refreshed using the refresh token provided
		strava.WithAutoRefresh(c.Context))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx := context.Background()
	act, err := client.Activity.Activity(ctx, 4273918760)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(act)
	os.Exit(0)
}
```
