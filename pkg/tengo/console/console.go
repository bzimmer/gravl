package console

import (
	"fmt"
	"io"
	"strings"

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

func replPrintln(out io.Writer) *tengo.UserFunction {
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

type Console struct {
	globals     []tengo.Object
	readline    *readline.Instance
	symbolTable *tengo.SymbolTable
}

func NewConsole(readline *readline.Instance) *Console {
	c := &Console{
		readline:    readline,
		symbolTable: tengo.NewSymbolTable(),
		globals:     make([]tengo.Object, tengo.GlobalsSize),
	}
	for idx, fn := range tengo.GetAllBuiltinFunctions() {
		c.symbolTable.DefineBuiltin(idx, fn.Name)
	}
	fn := replPrintln(readline.Config.Stdout)
	symbol := c.symbolTable.Define(fn.Name)
	c.globals[symbol.Index] = fn
	return c
}

func (c *Console) Run(modules *tengo.ModuleMap) error {
	acc := NewAccumulator()
	fileSet := parser.NewFileSet()
	var constants []tengo.Object
	for {
		line, err := c.readline.Readline()
		if err != nil {
			return err
		}
		acc.Push(line)

		cmd := acc.String()
		srcFile := fileSet.AddFile("repl", -1, len(cmd))
		p := parser.NewParser(srcFile, []byte(cmd), nil)
		file, err := p.ParseFile()
		if err != nil {
			list, ok := err.(parser.ErrorList)
			if ok && len(list) == 1 {
				x := list[0]
				if strings.Contains(x.Error(), "found 'EOF'") {
					c.readline.SetPrompt("... ")
					continue
				}
			}
			fmt.Println(err.Error())
			c.readline.SetPrompt(">>> ")
			acc.Reset()
			continue
		}

		c.readline.SetPrompt(">>> ")
		acc.Reset()

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
