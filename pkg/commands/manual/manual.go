package manual

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/analysis"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/commands/internal"
)

type command struct {
	Cmd     *cli.Command
	Lineage []*cli.Command
}

func (c *command) String() string { return c.fullname(" ") }

func (c *command) aliases() []string {
	if c.Cmd.Aliases == nil {
		return nil
	}
	var s []string
	for i := range c.Cmd.Aliases {
		if c.Cmd.Aliases[i] != "" {
			s = append(s, c.Cmd.Aliases[i])
		}
	}
	return s
}

func (c *command) fullname(sep string) string {
	var names []string
	for i := range c.Lineage {
		names = append(names, c.Lineage[i].Name)
	}
	names = append(names, c.Cmd.Name)
	return strings.Join(names, sep)
}

func ticks() string { return "```" }

func analyzersTemplate() (*template.Template, error) {
	return template.New("analyzers").
		Funcs(map[string]interface{}{
			"flags": func(s *flag.FlagSet) []*flag.Flag {
				var v []*flag.Flag
				if s != nil {
					s.VisitAll(func(f *flag.Flag) {
						v = append(v, f)
					})
				}
				sort.SliceStable(v, func(i, j int) bool {
					return v[i].Name < v[j].Name
				})
				return v
			},
			"ticks": ticks,
		}).
		Parse(`
{{- range .Analyzers }}

## *{{ .Name }}*

**Description**

{{ .Doc }}

{{- with .Flags }}

**Flags**

|Flag|Default|Description|
|-|-|-|
{{- range flags . }}
|{{ticks}}{{.Name}}{{ticks}}|{{ticks}}{{.DefValue}}{{ticks}}|{{.Usage}}|
{{- end }}
{{- end }}
{{- end }}
`)
}

func manualTemplate(root string) (*template.Template, error) {
	return template.New("commands").
		Funcs(map[string]interface{}{
			"usage": func(fn string) (string, error) {
				var err error
				path := filepath.Join(root, "docs", "commands", fn+".md")
				if _, err = os.Stat(path); os.IsNotExist(err) {
					// ok to skip any commands without usage documentation
					log.Warn().Str("path", path).Str("command", fn).Msg("missing")
					return "", nil
				}
				file, err := os.Open(path)
				if err != nil {
					return "", err
				}
				var usage []byte
				usage, err = ioutil.ReadAll(file)
				if err != nil {
					return "", err
				}
				log.Info().Str("path", path).Str("command", fn).Msg("reading")
				return strings.TrimSpace(string(usage)), nil
			},
			"join": func(s []string, sep string) string {
				return strings.Join(s, sep)
			},
			"fullname": func(c *command, sep string) string {
				return c.fullname(sep)
			},
			"aliases": func(c *command) []string {
				return c.aliases()
			},
			"names": func(f cli.Flag) string {
				// the first name is always the long name so skip it
				if len(f.Names()) <= 1 {
					return ""
				}
				return fmt.Sprintf("```%s```", strings.Join(f.Names()[1:], ", "))
			},
			"description": func(f cli.Flag) string {
				if x, ok := f.(cli.DocGenerationFlag); ok {
					return x.GetUsage()
				}
				return ""
			},
			"ticks": ticks,
		}).
		Parse(`# {{ .Name }} - {{ .Description }}

{{ usage "gravl" }}

## Commands

{{- range .Commands }}
* [{{ fullname . " " }}](#{{ fullname . "-" }})
{{- end }}

{{- range .Commands }}

## *{{ fullname . " " }}*

**Description**

{{ if .Cmd.Description }}{{ .Cmd.Description }}{{ else }}{{ .Cmd.Usage }}{{ end }}

{{ if .Cmd.Action }}
**Syntax**

{{ ticks }}sh
$ gravl {{ fullname . " " }} [flags] {{- if .Cmd.ArgsUsage }} {{.Cmd.ArgsUsage}}{{ end }}
{{ ticks }}
{{ end }}

{{- with .Cmd.Flags }}
**Flags**

|Name|Aliases|Description|
|-|-|-|
{{- range $f := . }}
|{{ticks}}{{ $f.Name }}{{ticks}}|{{ names $f }}|{{ description $f }}|
{{- end }}
{{- end }}

{{- $fn := fullname . "-" }}
{{- $x := usage $fn }}
{{ if $x }}
{{ if .Cmd.Action }}**Example**{{ else }}**Overview**{{ end }}

{{ $x }}
{{- end }}
{{- end }}
`)
}

func lineate(cmds []*cli.Command, lineage []*cli.Command) []*command {
	var commands []*command
	for i := range cmds {
		if cmds[i].Hidden {
			continue
		}
		commands = append(commands, &command{Cmd: cmds[i], Lineage: lineage})
		commands = append(commands, lineate(cmds[i].Subcommands, append(lineage, cmds[i]))...)
	}
	sort.SliceStable(commands, func(i, j int) bool {
		return commands[i].fullname("") < commands[j].fullname("")
	})
	return commands
}

var Manual = &cli.Command{
	Name:    "manual",
	Usage:   "Generate the `gravl` manual",
	Aliases: []string{"md"},
	Hidden:  true,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "commands",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "analyzers",
			Value: false,
		},
	},
	Before: func(c *cli.Context) error {
		if c.Bool("commands") == c.Bool("analyzers") {
			switch {
			case c.Bool("commands"):
				return errors.New("only one of `commands` or `analyzers` may be enabled at a time")
			case !c.Bool("commands"):
				return errors.New("one of `commands` or `analyzers` must be enabled")
			}
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		root, err := internal.Root()
		if err != nil {
			if os.IsNotExist(err) {
				return errors.New("manual creation can only happen from within the source tree")
			}
			return err
		}
		buffer := &bytes.Buffer{}
		switch {
		case c.Bool("commands"):
			// generate the command manual
			commands := lineate(c.App.Commands, nil)
			t := template.Must(manualTemplate(root))
			if err = t.Execute(buffer, map[string]interface{}{
				"Name":        c.App.Name,
				"Description": c.App.Description,
				"Commands":    commands,
			}); err != nil {
				return err
			}
		case c.Bool("analyzers"):
			// generate the analyzer manual
			a := analysis.All()
			sort.SliceStable(a, func(i, j int) bool {
				return a[i].Name < a[j].Name
			})
			t := template.Must(analyzersTemplate())
			if err = t.Execute(buffer, map[string]interface{}{
				"Analyzers": a,
			}); err != nil {
				return err
			}
		}
		fmt.Fprint(c.App.Writer, buffer.String())
		return nil
	},
}

var Commands = &cli.Command{
	Name:  "commands",
	Usage: "Return all possible commands",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "relative",
			Aliases: []string{"r"},
			Usage:   "Specify the command relative to the current working directory",
		},
	},
	Action: func(c *cli.Context) error {
		cmd := c.App.Name
		if c.Bool("relative") {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			cmd, err = os.Executable()
			if err != nil {
				return err
			}
			cmd, err = filepath.Rel(cwd, cmd)
			if err != nil {
				return err
			}
		}
		var s []string
		for _, c := range lineate(c.App.Commands, nil) {
			s = append(s, cmd+" "+c.fullname(" "))
		}
		return encoding.Encode(s)
	},
}
