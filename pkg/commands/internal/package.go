package internal

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-cmd/cmd"
)

func packageRoot() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// splitting the path drops the leading slash
	paths := make([]string, 1)
	paths[0] = "/"
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

// PackageGravl is primarily used for integration testing
func PackageGravl() string {
	root, err := packageRoot()
	if err != nil {
		panic(err)
	}
	return filepath.Join(root, "dist", "gravl")
}

// GravlCmd executes `gravl` commands
type GravlCmd struct {
	*cmd.Cmd
}

// Gravl creates a new instance of `Gravl`
func Gravl(args ...string) *GravlCmd {
	return &GravlCmd{Cmd: cmd.NewCmd(PackageGravl(), args...)}
}

// Success returns `true` if the `gravl` exit status is 0
func (c *GravlCmd) Success() bool {
	return c.Status().Exit == 0
}

// Stdout returns all the contents of stdout
func (c *GravlCmd) Stdout() string {
	return strings.Join(c.Status().Stdout, "\n")
}
