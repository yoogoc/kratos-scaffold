package wire

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"github.com/pkg/errors"
)

func AddToProviderSet(filePath string, constructorNames ...string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return errors.Wrapf(err, "parse file %s", filePath)
	}

	for _, name := range constructorNames {
		if err := addOneToProviderSet(f, name); err != nil {
			return err
		}
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return errors.Wrap(err, "format AST")
	}
	return os.WriteFile(filePath, buf.Bytes(), 0o644)
}

func addOneToProviderSet(f *ast.File, constructorName string) error {
	call := findNewSetCall(f)
	if call == nil {
		return errors.New("wire.NewSet call not found in file")
	}

	for _, arg := range call.Args {
		if ident, ok := arg.(*ast.Ident); ok && ident.Name == constructorName {
			return nil
		}
	}

	call.Args = append(call.Args, ast.NewIdent(constructorName))
	return nil
}

func findNewSetCall(f *ast.File) *ast.CallExpr {
	var result *ast.CallExpr
	ast.Inspect(f, func(n ast.Node) bool {
		if result != nil {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		if ident.Name == "wire" && sel.Sel.Name == "NewSet" {
			result = call
			return false
		}
		return true
	})
	return result
}
