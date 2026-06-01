package server

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
)

func RegisterService(serverDir, serviceName string, genHttp bool) error {
	if genHttp {
		if err := registerToHTTP(serverDir, serviceName); err != nil {
			return errors.Wrap(err, "register to http server")
		}
	} else {
		if err := registerToGRPC(serverDir, serviceName); err != nil {
			return errors.Wrap(err, "register to grpc server")
		}
	}
	return nil
}

func registerToGRPC(serverDir, serviceName string) error {
	filePath := path.Join(serverDir, "grpc.go")
	return registerToServer(filePath, "NewGRPCServer", serviceName, "ServiceServer")
}

func registerToHTTP(serverDir, serviceName string) error {
	filePath := path.Join(serverDir, "http.go")
	return registerToServer(filePath, "NewHTTPServer", serviceName, "ServiceHTTPServer")
}

func registerToServer(filePath, funcName, serviceName, registerSuffix string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return errors.Wrapf(err, "parse file %s", filePath)
	}

	funcDecl := findFuncDecl(f, funcName)
	if funcDecl == nil {
		return errors.Errorf("function %s not found in %s", funcName, filePath)
	}

	paramName := strcase.ToLowerCamel(serviceName)
	paramType := serviceName + "Service"

	if hasParam(funcDecl, paramName) {
		return nil
	}

	addParam(funcDecl, paramName, "service", paramType)

	registerCall := buildRegisterCall(serviceName, paramName, registerSuffix)
	insertBeforeReturn(funcDecl, registerCall)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return errors.Wrap(err, "format AST")
	}

	content, err := imports.Process(filePath, buf.Bytes(), nil)
	if err != nil {
		return errors.Wrap(err, "process imports")
	}

	return os.WriteFile(filePath, content, 0o644)
}

func findFuncDecl(f *ast.File, name string) *ast.FuncDecl {
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if ok && fd.Name.Name == name {
			return fd
		}
	}
	return nil
}

func hasParam(fd *ast.FuncDecl, paramName string) bool {
	for _, p := range fd.Type.Params.List {
		for _, n := range p.Names {
			if n.Name == paramName {
				return true
			}
		}
	}
	return false
}

func addParam(fd *ast.FuncDecl, paramName, pkgName, typeName string) {
	newParam := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(paramName)},
		Type: &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(pkgName),
				Sel: ast.NewIdent(typeName),
			},
		},
	}
	fd.Type.Params.List = append(fd.Type.Params.List, newParam)
}

// buildRegisterCall builds: v1.Register<Name>ServiceServer(srv, <paramName>)
// or for HTTP: v1.Register<Name>ServiceHTTPServer(srv, <paramName>)
func buildRegisterCall(serviceName, paramName, registerSuffix string) *ast.ExprStmt {
	// proto service name is "<Name>Service", so protoc generates:
	//   gRPC:  Register<Name>ServiceServer
	//   HTTP:  Register<Name>ServiceHTTPServer
	registerFuncName := "Register" + serviceName + registerSuffix
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("v1"),
				Sel: ast.NewIdent(registerFuncName),
			},
			Args: []ast.Expr{
				ast.NewIdent("srv"),
				ast.NewIdent(paramName),
			},
		},
	}
}

func insertBeforeReturn(fd *ast.FuncDecl, stmt ast.Stmt) {
	if fd.Body == nil {
		return
	}
	stmts := fd.Body.List
	for i := len(stmts) - 1; i >= 0; i-- {
		if _, ok := stmts[i].(*ast.ReturnStmt); ok {
			newStmts := make([]ast.Stmt, 0, len(stmts)+1)
			newStmts = append(newStmts, stmts[:i]...)
			newStmts = append(newStmts, stmt)
			newStmts = append(newStmts, stmts[i:]...)
			fd.Body.List = newStmts
			return
		}
	}
	fd.Body.List = append(fd.Body.List, stmt)
}
