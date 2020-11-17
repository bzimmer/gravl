package tengo

import (
	"fmt"
	"io"
	"strings"

	"github.com/bzimmer/gravl/pkg"
	"github.com/chzyer/readline"
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/parser"
)

func decorate(file *parser.File) *parser.File {
	var stmts []parser.Stmt
	for _, s := range file.Stmts {
		switch s := s.(type) {
		case *parser.ExprStmt:
			stmts = append(stmts, &parser.ExprStmt{
				Expr: &parser.CallExpr{
					Func: &parser.Ident{Name: "__repl_println__"},
					Args: []parser.Expr{s.Expr},
				},
			})
		// case *parser.AssignStmt:
		// 	stmts = append(stmts, s)
		// 	stmts = append(stmts, &parser.ExprStmt{
		// 		Expr: &parser.CallExpr{
		// 			Func: &parser.Ident{
		// 				Name: "__repl_println__",
		// 			},
		// 			Args: s.LHS,
		// 		},
		// 	})
		default:
			stmts = append(stmts, s)
		}
	}
	return &parser.File{
		InputFile: file.InputFile,
		Stmts:     stmts,
	}
}

func tengoPrintln(out io.Writer) *tengo.UserFunction {
	return &tengo.UserFunction{
		Name: "__repl_println__",
		Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
			var printArgs []interface{}
			for _, arg := range args {
				if _, isUndefined := arg.(*tengo.Undefined); !isUndefined {
					s, _ := tengo.ToString(arg)
					printArgs = append(printArgs, s+"\n")
				}
			}
			fmt.Fprint(out, printArgs...)
			return
		},
	}
}

func tengoVersion() *tengo.UserFunction {
	return &tengo.UserFunction{
		Name: "version",
		Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
			x, err := tengo.FromInterface(pkg.BuildVersion)
			if err != nil {
				return nil, err
			}
			return x, nil
		},
	}
}

type accumulator struct {
	lines []string
}

func (a *accumulator) String() string {
	if len(a.lines) == 0 {
		return ""
	}
	return strings.Join(a.lines, "")
}

func (a *accumulator) push(line string) {
	a.lines = append(a.lines, line)
}

func (a *accumulator) reset() *accumulator {
	a.lines = make([]string, 0)
	return a
}

type Console struct {
	readline    *readline.Instance
	symbolTable *tengo.SymbolTable
	globals     []tengo.Object
}

func NewConsole(readline *readline.Instance) *Console {
	return &Console{
		readline: readline,
	}
}

func (c *Console) init() {
	c.symbolTable = tengo.NewSymbolTable()
	c.globals = make([]tengo.Object, tengo.GlobalsSize)

	for idx, fn := range tengo.GetAllBuiltinFunctions() {
		c.symbolTable.DefineBuiltin(idx, fn.Name)
	}

	for _, f := range []*tengo.UserFunction{tengoPrintln(c.readline.Config.Stdout), tengoVersion()} {
		symbol := c.symbolTable.Define(f.Name)
		c.globals[symbol.Index] = f
	}
}

func (c *Console) Run(modules *tengo.ModuleMap) error {
	c.init()

	acc := (&accumulator{}).reset()
	fileSet := parser.NewFileSet()
	var constants []tengo.Object
	for {
		line, err := c.readline.Readline()
		if err != nil {
			return err
		}
		acc.push(line)

		cmd := acc.String()
		srcFile := fileSet.AddFile("repl", -1, len(cmd))
		p := parser.NewParser(srcFile, []byte(cmd), nil)
		file, err := p.ParseFile()
		if err != nil {
			list, ok := err.(parser.ErrorList)
			if ok {
				x := list[0]
				if strings.Contains(x.Error(), "found 'EOF'") {
					c.readline.SetPrompt("... ")
				}
			}
			continue
		}

		c.readline.SetPrompt(">>> ")
		acc.reset()

		file = decorate(file)
		compiler := tengo.NewCompiler(srcFile, c.symbolTable, constants, modules, nil)
		if err := compiler.Compile(file); err != nil {
			fmt.Fprint(c.readline.Config.Stdout, err.Error())
			continue
		}

		bytecode := compiler.Bytecode()
		machine := tengo.NewVM(bytecode, c.globals, -1)
		if err := machine.Run(); err != nil {
			fmt.Fprint(c.readline.Config.Stdout, err.Error())
			continue
		}
		constants = bytecode.Constants
	}
}
