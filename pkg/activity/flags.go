package activity

import (
	"time"

	"github.com/urfave/cli/v2"
)

var RateLimitFlags = []cli.Flag{
	&cli.DurationFlag{
		Name:  "rate-limit",
		Value: time.Millisecond * 1500,
		Usage: "Minimum time interval between API request events (eg, 1ms, 2s, 5m, 3h)",
	},
	&cli.IntFlag{
		Name:  "rate-burst",
		Value: 35,
		Usage: "Maximum burst size for API request events",
	},
	&cli.IntFlag{
		Name:  "concurrency",
		Value: 2,
		Usage: "Allowable concurrent queries to data providers",
	},
}
