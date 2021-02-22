package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-cmd/cmd"
	"github.com/rs/zerolog/log"
)

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
		root := filepath.Join(x, "go.mod")
		if _, err := os.Stat(root); os.IsNotExist(err) {
			paths = paths[:len(paths)-1]
		} else {
			return x, nil
		}
	}
	return "", os.ErrNotExist
}

// GravlCmd executes `gravl` commands
type GravlCmd struct {
	*cmd.Cmd
}

// Gravl creates a new instance of `Gravl`
func Gravl(args ...string) *GravlCmd {
	root, err := root()
	if err != nil {
		panic(err)
	}
	g := filepath.Join(root, "dist", "gravl")
	log.Info().Strs("args", args).Msg("cmd")
	return &GravlCmd{Cmd: cmd.NewCmd(g, args...)}
}

// Success returns `true` if the `gravl` exit status is 0
func (c *GravlCmd) Success() bool {
	b := c.Status().Exit == 0
	if !b {
		fmt.Fprint(os.Stderr, c.Stderr())
		fmt.Fprint(os.Stderr, c.Stdout())
	}
	return b
}

// Stdout returns all the contents of stdout
func (c *GravlCmd) Stdout() string {
	return strings.Join(c.Status().Stdout, "\n")
}

// Stderr returns all the contents of stderr
func (c *GravlCmd) Stderr() string {
	return strings.Join(c.Status().Stderr, "\n")
}
