package manual

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/bzimmer/gravl/pkg"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
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

func read(path string) (string, error) {
	var err error
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	var usage []byte
	usage, err = ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(usage), nil
}

func manualTemplate(root string) (*template.Template, error) {
	man, err := read(filepath.Join(root, "docs", "commands", "_commands.md"))
	if err != nil {
		return nil, err
	}
	return template.New("commands").
		Funcs(map[string]interface{}{
			"partial": func(fn string) (string, error) {
				path := filepath.Join(root, "docs", "commands", fn+".md")
				usage, err := read(path)
				if err != nil {
					if os.IsNotExist(err) {
						// ok to skip any commands without usage documentation
						log.Warn().Str("path", path).Str("command", fn).Msg("missing")
						return "", nil
					}
					return "", err
				}
				log.Info().Str("path", path).Str("command", fn).Msg("reading")
				return usage, nil
			},
			"join": strings.Join,
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
				return strings.Join(f.Names()[1:], ", ")
			},
			"description": func(f cli.Flag) string {
				if x, ok := f.(cli.DocGenerationFlag); ok {
					return x.GetUsage()
				}
				return ""
			},
		}).
		Parse(man)
}

func lineate(cmds, lineage []*cli.Command) []*command {
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

func Command() *cli.Command {
	return &cli.Command{
		Name:    "manual",
		Usage:   "Generate the `gravl` manual",
		Aliases: []string{"md"},
		Hidden:  true,
		Action: func(c *cli.Context) error {
			var buffer bytes.Buffer
			commands := lineate(c.App.Commands, nil)
			t, err := manualTemplate(".")
			if err != nil {
				return err
			}
			if err := t.Execute(&buffer, map[string]interface{}{
				"Name":        c.App.Name,
				"Description": c.App.Description,
				"GlobalFlags": c.App.Flags,
				"Commands":    commands,
			}); err != nil {
				return err
			}
			fmt.Fprint(c.App.Writer, buffer.String())
			return nil
		},
	}
}

func Commands() *cli.Command {
	return &cli.Command{
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
				log.Info().Str("name", c.Cmd.Name).Str("usage", c.Cmd.Usage).Msg("command")
				s = append(s, cmd+" "+c.fullname(" "))
			}
			return pkg.Runtime(c).Encoder.Encode(s)
		},
	}
}
