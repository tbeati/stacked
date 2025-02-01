package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/beati/stacked"
	"github.com/beati/stacked/cmd/stacked/analyser"
)

var src = `
package foo

import "os"

func assign() error {
	f, err := os.Create("test")
	if err != nil {
		return err
	}

	n, err := f.WriteString("test")
	if err != nil {
		return err
	}
	_ = n

	err = f.Close()
	if err != nil {
		return err
	}

	err = os.Chdir("test")
	if err != nil {
		return err
	}

	name, err := os.Hostname()
	if err != nil {
		return err
	}
	_ = name

	name, err = "test", os.Chdir("test")
	if err != nil {
		return err
	}

	return nil
}
`

func main() {
	singlechecker.Main(analyzer.Analyzer)
}

func addWrap() error {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, "", src, parser.AllErrors)
	if err != nil {
		return stacked.Wrap(err)
	}

	conf := types.Config{
		Importer: importer.Default(),
	}
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}
	_, err = conf.Check("", fs, []*ast.File{node}, info)
	if err != nil {
		return stacked.Wrap(err)
	}

	//fmt.Println(info)

	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		switch stmt := n.(type) {
		case *ast.AssignStmt:
			fmt.Printf("AssignStmt: %#v\n", stmt)
			handleAssignment(stmt.Lhs, stmt.Rhs)
		}

		return true
	})

	/*
		ast.Inspect(node, func(n ast.Node) bool {
			if n == nil {
				return true
			}

			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			fmt.Printf("CallExpr: %v\n", callExpr.Fun)

			tv, ok := info.Types[callExpr.Fun]
			if !ok {
				return true
			}

			if returnsErrorType(tv.Type) {
				fmt.Printf("Function call: %s (returns error)\n", callExpr.Fun)
			}

			return true
		})
	*/

	return nil
}

func handleAssignment(lhs, rhs []ast.Expr) {
	if len(lhs) == len(rhs) {
	} else {
		//call, ok := rhs[0].(*ast.CallExpr)
		//if ok {
		//}
	}
}

func returnsErrorType(typ types.Type) bool {
	if sig, ok := typ.(*types.Signature); ok {
		results := sig.Results()
		for i := 0; i < results.Len(); i++ {
			if results.At(i).Type().String() == "error" {
				return true
			}
		}
	}
	return false
}
