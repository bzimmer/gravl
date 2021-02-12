package manual

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
)

var manualTemplate = `
# NAME

{{ .App.Name }}{{ if .App.Usage }} - {{ .App.Usage }}{{ end }}

# SYNOPSIS

{{ .App.Name }}
`

// {{ if .SynopsisArgs }}
// ` + "```" + `
// {{ range $v := .SynopsisArgs }}{{ $v }}{{ end }}` + "```" + `
// {{ end }}{{ if .App.UsageText }}
// # DESCRIPTION

// {{ .App.UsageText }}
// {{ end }}
// **Usage**:

// ` + "```" + `
// {{ .App.Name }} [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
// ` + "```" + `
// {{ if .GlobalArgs }}
// # GLOBAL OPTIONS
// {{ range $v := .GlobalArgs }}
// {{ $v }}{{ end }}
// {{ end }}{{ if .Commands }}
// # COMMANDS
// {{ range $v := .Commands }}
// {{ $v }}{{ end }}{{ end }}`

type templateEnv struct {
	App        *cli.App
	Commands   []*cli.Command
	GlobalArgs []cli.Flag
}

func manual(c *cli.Context) error {
	const name = "manual"
	t, err := template.New(name).Parse(manualTemplate)
	if err != nil {
		return err
	}
	app := c.App
	return t.ExecuteTemplate(app.Writer, name, &templateEnv{
		App:        app,
		Commands:   app.Commands,
		GlobalArgs: app.VisibleFlags(),
	})
}

var markdownCommand = &cli.Command{
	Name:    "markdown",
	Usage:   "Generate the manual in Markdown",
	Aliases: []string{"md"},
	Action:  manual,
}

var commandsCommand = &cli.Command{
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
		fmt.Printf("%v\n", c.App.Commands)
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

var Command = &cli.Command{
	Name:     "manual",
	Category: "manual",
	Usage:    "Generate a manual for the cli",
	Subcommands: []*cli.Command{
		commandsCommand,
		markdownCommand,
	},
}
