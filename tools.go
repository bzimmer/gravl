// +build tools

package main

import (
	// brew install go-task
	// _ "github.com/go-task/task/v3/cmd/task"
	// brew install golangci-lint
	// _ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/posener/goreadme/cmd/goreadme"
	_ "github.com/securego/gosec/v2/cmd/gosec"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/tools/cmd/stringer"
)
