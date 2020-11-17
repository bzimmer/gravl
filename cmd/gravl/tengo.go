package gravl

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/tengo/console"
	gravllib "github.com/bzimmer/gravl/pkg/tengo/stdlib"
)

func modules() *tengo.ModuleMap {
	m := tengo.NewModuleMap()
	m.AddMap(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	m.AddMap(gravllib.GetModuleMap(gravllib.AllModuleNames()...))
	return m
}

var tengoCommand = &cli.Command{
	Name:     "tengo",
	Category: "api",
	Usage:    "Run tengo",
	Action: func(c *cli.Context) error {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		rl, err := readline.NewEx(&readline.Config{
			Prompt:            ">>> ",
			HistoryFile:       filepath.Join(home, ".gravl_history"),
			InterruptPrompt:   "^C",
			EOFPrompt:         "exit",
			Stdin:             ioutil.NopCloser(c.App.Reader),
			Stdout:            c.App.Writer,
			HistorySearchFold: true,
		})
		if err != nil {
			return err
		}
		console := console.NewConsole(rl)
		err = console.Run(modules())
		if err == io.EOF {
			return nil
		}
		return err
	},
}
