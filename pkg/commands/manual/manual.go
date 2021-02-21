package manual

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/analysis"
	"github.com/bzimmer/gravl/pkg/commands/encoding"
	"github.com/bzimmer/gravl/pkg/commands/internal"
)

type command struct {
	Cmd     *cli.Command
	Lineage []*cli.Command
}

func ticks() string { return "```" }

func (c *command) fullname(sep string) string {
	var names []string
	for i := range c.Lineage {
		names = append(names, c.Lineage[i].Name)
	}
	names = append(names, c.Cmd.Name)
	return strings.Join(names, sep)
}

func analysisTemplate(root string) (*template.Template, error) {
	if root == "" {
		return nil, nil
	}
	return template.New("analysis").
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
`)
}

func commandTemplate(root string) (*template.Template, error) {
	return template.New("command").
		Funcs(map[string]interface{}{
			"usage": func(c *command) (string, error) {
				var err error
				fn := filepath.Join(root, "docs", "usage", c.fullname("-")+".md")
				if _, err = os.Stat(fn); os.IsNotExist(err) {
					// ok to skip any commands without usage documentation
					return "", nil
				}
				file, err := os.Open(fn)
				if err != nil {
					return "", err
				}
				var usage []byte
				usage, err = ioutil.ReadAll(file)
				if err != nil {
					return "", err
				}
				return strings.TrimSpace(string(usage)), nil
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
			"lineage": func(c *command) string {
				return c.fullname(" ")
			},
			"ticks": ticks,
		}).
		Parse(`
{{- if .Cmd.Action }}
## *{{ lineage . }}*

**Description**

{{ if .Cmd.Description }}{{ .Cmd.Description }}{{ else }}{{ .Cmd.Usage }}{{ end }}

**Syntax:**

{{ ticks }}sh
$ gravl {{ lineage . }}{{- if .Cmd.ArgsUsage }} {{.Cmd.ArgsUsage}}{{ end }}
{{ ticks }}

{{- with .Cmd.Flags }}
**Flags:**

|Flag|Short|Description|
|-|-|-|
{{- range $f := . }}
|{{ticks}}{{ $f.Name }}{{ticks}}|{{ names $f }}|{{ description $f }}|
{{- end }}
{{ end }}

{{- with $x := usage . }}
**Example:**

{{ . }}
{{- end }}
{{- end }}
`)
}

func manual(t *template.Template, buffer io.Writer, cmds []*cli.Command, lineage []*cli.Command) error {
	for i := range cmds {
		c := &command{Cmd: cmds[i], Lineage: lineage}
		if err := t.Execute(buffer, c); err != nil {
			return err
		}
		if err := manual(t, buffer, cmds[i].Subcommands, append(lineage, cmds[i])); err != nil {
			return err
		}
	}
	return nil
}

func analyzer(t *template.Template, buffer io.Writer) error {
	a := analysis.All()
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].Name < a[j].Name
	})
	for i := range a {
		if err := t.Execute(buffer, a[i]); err != nil {
			return err
		}
	}
	return nil
}

var Manual = &cli.Command{
	Name:    "manual",
	Usage:   "Generate the `gravl` manual",
	Aliases: []string{"md"},
	Hidden:  true,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "manual",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "analyzer",
			Value: false,
		},
	},
	Before: func(c *cli.Context) error {
		switch {
		case c.Bool("manual") && (c.Bool("manual") == c.Bool("analyzer")):
			return errors.New("only one of `manual` or `analyzer` may be enabled at a time")
		case !c.Bool("manual") && (c.Bool("manual") == c.Bool("analyzer")):
			return errors.New("one of `manual` or `analyzer` must be enabled")
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
		switch {
		case c.Bool("manual"):
			// generate the main manual
			buffer := &bytes.Buffer{}
			t, err := template.New("manual").
				Parse(`
# {{ .Name }} - {{ .Description }}
`)
			if err != nil {
				return err
			}
			if err = t.Execute(buffer, c.App); err != nil {
				return err
			}
			t, err = commandTemplate(root)
			if err != nil {
				return err
			}
			if err = manual(t, buffer, c.App.Commands, nil); err != nil {
				return err
			}
			fmt.Fprint(c.App.Writer, buffer.String())
		case c.Bool("analyzer"):
			// generate the analyzer manual
			buffer := &bytes.Buffer{}
			t, err := analysisTemplate(root)
			if err != nil {
				return err
			}
			if err = analyzer(t, buffer); err != nil {
				return nil
			}
			fmt.Fprint(c.App.Writer, buffer.String())
		}
		return nil
	},
}

var Commands = &cli.Command{
	Name:   "commands",
	Usage:  "Return all possible commands",
	Hidden: true,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "relative",
			Aliases: []string{"r"},
			Usage:   "Specify the command relative to the current working directory",
		},
	},
	Action: func(c *cli.Context) error {
		var commander func(string, []*cli.Command) []string
		commander = func(prefix string, cmds []*cli.Command) []string {
			var commands []string
			for i := range cmds {
				cmd := fmt.Sprintf("%s %s", prefix, cmds[i].Name)
				if !cmds[i].Hidden && cmds[i].Action != nil {
					commands = append(commands, cmd)
				}
				commands = append(commands, commander(cmd, cmds[i].Subcommands)...)
			}
			return commands
		}
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
		commands := commander(cmd, c.App.Commands)
		return encoding.Encode(commands)
	},
}
