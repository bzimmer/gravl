package activity

import (
	"errors"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/rs/zerolog/log"
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

type DateParserFunc func(string) (time.Time, error)

func NaturalParse(pattern string) (time.Time, error) {
	return naturaldate.Parse(pattern, time.Now())
}

func AraddonParse(pattern string) (time.Time, error) {
	return dateparse.ParseStrict(pattern)
}

func parse(pattern string, parsers []DateParserFunc) (time.Time, error) {
	for i := range parsers {
		date, err := parsers[i](pattern)
		if err == nil {
			return date.In(time.UTC), nil
		}
		log.Debug().Err(err).Msg("parse")
	}
	return time.Time{}, fmt.Errorf("failed to parse: %s", pattern)
}

// DateRange returns the date range specified in the command line flags
func DateRange(c *cli.Context, parsers ...DateParserFunc) (time.Time, time.Time, error) {
	var err error
	var before, after time.Time
	if len(parsers) == 0 {
		return DateRange(c, NaturalParse)
	}
	if c.IsSet("before") {
		before, err = parse(c.String("before"), parsers)
		if err != nil {
			before, after = time.Time{}, time.Time{}
			return before, after, err
		}
	}
	if c.IsSet("after") {
		if before.IsZero() {
			before = time.Now()
		}
		after, err = parse(c.String("after"), parsers)
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
