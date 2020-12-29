package analysis

//go:generate stringer -type=Units -linecomment -output=units_string.go

import "fmt"

type Units int

const (
	Imperial Units = iota // imperial
	Metric                // metric
)

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
