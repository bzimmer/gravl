package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang/gddo/gosrc"
	"github.com/posener/goreadme"
)

func gen() error {
	ctx := context.Background()
	gr := goreadme.New(&http.Client{}).
		WithConfig(goreadme.Config{
			Types:           false,
			Functions:       false,
			SkipSubPackages: true,
		})
	path, err := filepath.Abs("./")
	if err != nil {
		return err
	}
	gosrc.SetLocalDevMode(path)

	out, err := os.Create("README.md")
	if err != nil {
		return err
	}
	defer out.Close()
	return gr.Create(ctx, ".", out)
}

func main() {
	if err := gen(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
