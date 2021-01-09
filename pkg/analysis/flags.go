package analysis

//go:generate stringer -type=Units -linecomment -output=flags_string.go

import (
	"fmt"
	"time"
)

type Units int

const (
	Imperial Units = iota // imperial
	Metric                // metric
)

const YMD = "2006-01-02"

// UnitsFlag implements the flag.Value interface for analysis.Units
type UnitsFlag struct {
	Units Units
}

func (u *UnitsFlag) String() string {
	return u.Units.String()
}

func (u *UnitsFlag) Set(value string) error {
	switch value {
	case "imperial":
		u.Units = Imperial
	case "metric":
		u.Units = Metric
	default:
		return fmt.Errorf("unexpected unit '%s'", value)
	}
	return nil
}

func (u *UnitsFlag) Get() interface{} {
	return u.Units
}

// TimeFlag implements the flag.Value interface for time.Time
type TimeFlag struct {
	// The format of the string to parse
	Format string

	// The time value parsed from the flag.
	Time time.Time
}

func (t *TimeFlag) String() string {
	return format(t.Time, t.Format)
}

func (t *TimeFlag) Set(s string) error {
	var err error
	t.Time, err = parse(s, t.Format)
	return err
}

func (t *TimeFlag) Get() interface{} {
	return t.Time
}

func parse(s string, format string) (time.Time, error) {
	if format == "" {
		format = YMD
	}
	return time.Parse(format, s)
}

func format(t time.Time, format string) string {
	if format == "" {
		return fmt.Sprintf("%q", t.Format(YMD))
	}
	return fmt.Sprintf("%q", t.Format(format))
}
