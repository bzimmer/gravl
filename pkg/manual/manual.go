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

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg"
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

// root finds the root of the source tree by recursively ascending until 'go.mod' is located
func root() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	paths := []string{string(os.PathSeparator)}
	paths = append(paths, strings.Split(path, string(os.PathSeparator))...)
	for len(paths) > 0 {
		x := filepath.Join(paths...)
		mod := filepath.Join(x, "go.mod")
		if _, err := os.Stat(mod); os.IsNotExist(err) {
			paths = paths[:len(paths)-1]
		} else {
			return x, nil
		}
	}
	return "", os.ErrNotExist
}

func Command() *cli.Command {
	return &cli.Command{
		Name:    "manual",
		Usage:   "Generate the `gravl` manual",
		Aliases: []string{"man"},
		Hidden:  true,
		Action: func(c *cli.Context) error {
			var buffer bytes.Buffer
			commands := lineate(c.App.Commands, nil)
			path, err := root()
			if err != nil {
				return err
			}
			t, err := manualTemplate(path)
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
				cmd, err = filepath.Abs(cmd)
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

func Vars() *cli.Command {
	return &cli.Command{
		Name:   "vars",
		Hidden: true,
		Action: func(c *cli.Context) error {
			var vars []string
			for _, cmd := range lineate(c.App.Commands, nil) {
				for _, flag := range cmd.Cmd.Flags {
					switch v := flag.(type) {
					case *cli.StringFlag:
						vars = append(vars, v.EnvVars...)
					case *cli.BoolFlag:
						vars = append(vars, v.EnvVars...)
					case *cli.IntFlag:
						vars = append(vars, v.EnvVars...)
					}
				}
			}
			kv := make(map[string]bool)
			for _, env := range vars {
				kv[env] = true
			}
			var k []string
			for v := range kv {
				k = append(k, v)
			}
			sort.Strings(k)
			for _, v := range k {
				fmt.Fprintln(c.App.Writer, v+"=")
			}
			return nil
		},
	}
}
