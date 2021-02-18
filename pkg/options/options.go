package options

import (
	"encoding/csv"
	"errors"
	"flag"
	"io"
	"strings"
)

// Option represents a flag name and optional options
// For example:
//
//  $ gravl pass -a cluster,clusters=5
//
// would be represented by the Name 'cluster' and Options `clusters=5`
type Option struct {
	// Name is the flag name
	Name string
	// Options is the mapping of an option to a value
	Options map[string]string
}

// Parse the option string into an Option instance
func Parse(option string) (*Option, error) {
	reader := csv.NewReader(strings.NewReader(option))
	opts, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("missing options")
		}
		return nil, err
	}
	t := &Option{Name: opts[0], Options: make(map[string]string)}
	for i := 1; i < len(opts); i++ {
		x := strings.SplitN(opts[i], "=", 2)
		if len(x) != 2 {
			return nil, errors.New("missing '=' separating key from value (eg x=y)")
		}
		t.Options[x[0]] = x[1]
	}
	return t, nil
}

// ApplyFlags from the options to the flagset
func (t *Option) ApplyFlags(fs *flag.FlagSet) error {
	if fs == nil {
		return nil
	}
	var f []string
	for k, v := range t.Options {
		switch v {
		case "", "true", "false":
			f = append(f, "--"+k)
		default:
			f = append(f, "--"+k, v)
		}
	}
	return fs.Parse(f)
}
