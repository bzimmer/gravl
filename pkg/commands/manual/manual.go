package manual

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
)

type command struct {
	Cmd     *cli.Command
	Lineage []*cli.Command
}

func (c *command) fullname(sep string) string {
	var names []string
	for i := range c.Lineage {
		names = append(names, c.Lineage[i].Name)
	}
	names = append(names, c.Cmd.Name)
	return strings.Join(names, sep)
}

var tmpl = template.Must(template.New("command").
	Funcs(map[string]interface{}{
		"usage": func(c *command) string {
			s := usages[c.fullname("-")]
			usage, err := hex.DecodeString(s)
			if err != nil {
				log.Warn().Err(err).Msg("hex decode")
				return ""
			}
			return strings.TrimSpace(string(usage))
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
		"ticks": func() string { return "```" },
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
`))

func manual(buffer io.Writer, cmds []*cli.Command, lineage []*cli.Command) error {
	for i := range cmds {
		c := &command{Cmd: cmds[i], Lineage: lineage}
		if err := tmpl.Execute(buffer, c); err != nil {
			return err
		}
		if err := manual(buffer, cmds[i].Subcommands, append(lineage, cmds[i])); err != nil {
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
	Action: func(c *cli.Context) error {
		buffer := &bytes.Buffer{}
		t := template.Must(template.New("manual").
			Parse(`
# {{ .Name }} - {{ .Description }}
`))
		if err := t.Execute(buffer, c.App); err != nil {
			return err
		}
		if err := manual(buffer, c.App.Commands, nil); err != nil {
			return err
		}
		fmt.Fprintln(c.App.Writer, buffer.String())
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
