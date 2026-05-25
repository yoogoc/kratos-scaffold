package wire

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"

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

	output := formatNewSetMultiLine(buf.String())
	return os.WriteFile(filePath, []byte(output), 0o644)
}

var newSetRe = regexp.MustCompile(`wire\.NewSet\(([^)]*)\)`)

func formatNewSetMultiLine(src string) string {
	return newSetRe.ReplaceAllStringFunc(src, func(match string) string {
		inner := newSetRe.FindStringSubmatch(match)[1]
		inner = strings.TrimSpace(inner)
		if inner == "" {
			return match
		}
		args := strings.Split(inner, ",")
		for i := range args {
			args[i] = strings.TrimSpace(args[i])
		}
		// filter empty
		var cleaned []string
		for _, a := range args {
			if a != "" {
				cleaned = append(cleaned, a)
			}
		}
		if len(cleaned) == 0 {
			return match
		}
		var b strings.Builder
		b.WriteString("wire.NewSet(\n")
		for _, a := range cleaned {
			b.WriteString("\t")
			b.WriteString(a)
			b.WriteString(",\n")
		}
		b.WriteString(")")
		return b.String()
	})
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
