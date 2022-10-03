package activity

import (
	"errors"
	"time"

	"github.com/tj/go-naturaldate"
	"github.com/urfave/cli/v2"
)

// RateLimitFlags support specifying a rate limit for a query
func RateLimitFlags() []cli.Flag {
	return []cli.Flag{
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
			Usage: "Maximum concurrent API queries",
		},
	}
}

// DateRangeFlags support specifying a date range for a query
func DateRangeFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "after",
			Aliases: []string{"since"},
			Usage:   "Return results after the time specified",
		},
		&cli.StringFlag{
			Name:  "before",
			Usage: "Return results before the time specified",
		},
	}
}

// DateRange returns the date range specified in the command line flags
func DateRange(c *cli.Context) (time.Time, time.Time, error) {
	var err error
	var before, after time.Time
	if c.IsSet("before") {
		before, err = naturaldate.Parse(c.String("before"), time.Now())
		if err != nil {
			before, after = time.Time{}, time.Time{}
			return before, after, err
		}
	}
	if c.IsSet("after") {
		if before.IsZero() {
			before = time.Now()
		}
		after, err = naturaldate.Parse(c.String("after"), time.Now())
		if err != nil {
			before, after = time.Time{}, time.Time{}
			return before, after, err
		}
		if after.After(before) {
			before, after = time.Time{}, time.Time{}
			return before, after, errors.New("invalid date range")
		}
	}
	return before, after, nil
}
