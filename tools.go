// +build tools

package main

import (
	_ "github.com/go-task/task/v3/cmd/task"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/stringer"
)
