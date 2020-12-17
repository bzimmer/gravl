package main

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/fatih/structtag"
)

type Field struct {
	Name string
	Type string
	Base string
}

type Struct struct {
	Name   string
	Fields []*Field
}

func (s *Struct) add(f *Field) {
	if s.Fields == nil {
		s.Fields = make([]*Field, 0, 5)
	}
	s.Fields = append(s.Fields, f)
}

type Units struct {
	Package string
	Structs []*Struct
}

func (u *Units) add(s *Struct) {
	if u.Structs == nil {
		u.Structs = make([]*Struct, 0, 5)
	}
	u.Structs = append(u.Structs, s)
}

func parseFile(filename string) (*ast.File, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	return parser.ParseFile(token.NewFileSet(), filename, nil, 0)
}

func (u *Units) visit(n ast.Node) bool {
	var s *Struct
	switch x := n.(type) {
	case *ast.TypeSpec:
		switch x.Type.(type) {
		case *ast.StructType:
			s = &Struct{Name: x.Name.Name}
			st := x.Type.(*ast.StructType)
			for _, field := range st.Fields.List {
				if field.Tag == nil {
					continue
				}
				name := field.Names[0].Name
				tag := strings.ReplaceAll(field.Tag.Value, "`", "")
				tags, err := structtag.Parse(tag)
				if err != nil {
					log.Error().Err(err).Str("tag", tag).Msg("skipping")
					continue
				}
				units, err := tags.Get("units")
				if err != nil {
					continue
				}
				f := &Field{
					Name: name,
					Type: "",
					Base: units.Name,
				}
				s.add(f)
			}
			if len(s.Fields) > 0 {
				sort.Slice(s.Fields, func(i, j int) bool {
					return s.Fields[i].Name < s.Fields[j].Name
				})
				u.add(s)
			}
			s = nil
		}
	}
	sort.Slice(u.Structs, func(i, j int) bool {
		return u.Structs[i].Name < u.Structs[j].Name
	})
	return true
}

func visitUnits(f *ast.File) {
	u := &Units{
		Package: f.Name.String(),
		Structs: make([]*Struct, 0, 10),
	}
	ast.Inspect(f, u.visit)
	fmt.Printf("Package: %s\n\n", u.Package)
	for _, s := range u.Structs {
		fmt.Printf(" %s\n", s.Name)
		for _, f := range s.Fields {
			fmt.Printf("  %-30s %s\n", f.Name, f.Base)
		}
	}
}

func main() {
	app := &cli.App{
		Name:     "genunits",
		HelpName: "genunits",
		Before: func(c *cli.Context) error {
			n := c.NArg()
			if n != 1 {
				return errors.New("expected only one argument")
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			args := c.Args().Slice()
			f, err := parseFile(args[0])
			if err != nil {
				return err
			}
			visitUnits(f)
			return nil
		},
	}
	ctx := context.Background()
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
	os.Exit(0)
}
