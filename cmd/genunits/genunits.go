package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/structtag"
	"github.com/urfave/cli/v2"
)

// Field represents a single field with a unit
type Field struct {
	Name string
	Type string
	Base string
}

// Struct represents a struct containing fields with units
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

// Units contains all the structs in a package with fields with units
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

func parseType(t ast.Expr) string {
	for {
		switch s := t.(type) {
		case *ast.ArrayType:
			//  Type: (*ast.ArrayType)(0xc000073ec0)({
			//   Lbrack: (token.Pos) 1647,
			//   Len: (ast.Expr) <nil>,
			//   Elt: (*ast.SelectorExpr)(0xc0001a2940)({
			//    X: (*ast.Ident)(0xc0001a2900)(unit),
			//    Sel: (*ast.Ident)(0xc0001a2920)(Length)
			//   })
			t = s.Elt
		case *ast.SelectorExpr:
			// 	Type: (*ast.SelectorExpr)(0xc0000e4580)({
			//   X: (*ast.Ident)(0xc0000e4540)(unit),
			//   Sel: (*ast.Ident)(0xc0000e4560)(Length)
			//  })
			if s.X != nil {
				return fmt.Sprintf("%s.%s", s.X, s.Sel)
			}
			return s.Sel.String()
		default:
			return ""
		}
	}
}

func parseUnits(val string) (*structtag.Tag, error) {
	tag := strings.ReplaceAll(val, "`", "")
	tags, err := structtag.Parse(tag)
	if err != nil {
		return nil, err
	}
	return tags.Get("units")
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
				units, err := parseUnits(field.Tag.Value)
				if err != nil {
					if err.Error() == "tag does not exist" {
						continue
					}
					fmt.Fprintln(os.Stderr, err)
					return false
				}
				typ := parseType(field.Type)
				if typ == "" {
					spew.Dump(field)
					return false
				}
				f := &Field{
					Name: name,
					Type: typ,
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

func visitUnits(f *ast.File) *Units {
	u := &Units{
		Package: f.Name.String(),
		Structs: make([]*Struct, 0),
	}
	ast.Inspect(f, u.visit)
	return u
}

func collect(paths ...string) ([]string, error) {
	var files []string
	for _, path := range paths {
		err := filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if strings.HasSuffix(name, ".go") {
				files = append(files, name)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func main() {
	app := &cli.App{
		Name:     "genunits",
		HelpName: "genunits",
		Action: func(c *cli.Context) error {
			paths, err := collect(c.Args().Slice()...)
			if err != nil {
				return err
			}
			for _, path := range paths {
				f, err := parseFile(path)
				if err != nil {
					return err
				}
				u := visitUnits(f)
				if len(u.Structs) == 0 {
					continue
				}
				j := json.NewEncoder(c.App.Writer)
				if err := j.Encode(u); err != nil {
					return err
				}
			}
			return nil
		},
	}
	ctx := context.Background()
	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
