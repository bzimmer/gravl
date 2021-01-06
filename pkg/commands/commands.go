package commands

import "github.com/urfave/cli/v2"

func Merge(flags ...[]cli.Flag) []cli.Flag {
	var f []cli.Flag
	for _, x := range flags {
		f = append(f, x...)
	}
	return f
}

func Before(befores ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		for _, fn := range befores {
			if fn == nil {
				continue
			}
			if e := fn(c); e != nil {
				return e
			}
		}
		return nil
	}
}
