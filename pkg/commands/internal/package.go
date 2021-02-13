package internal

import (
	"os"
	"path/filepath"
	"strings"
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
