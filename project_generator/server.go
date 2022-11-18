package project_generator

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

type serverTmpl struct {
	AppPkgPath string
	ServerPath string
	GenGrpc    bool
	GenHttp    bool
}

type ServerTmplOption func(*serverTmpl)

func NonGenGrpc() ServerTmplOption {
	return func(st *serverTmpl) {
		st.GenGrpc = false
	}
}

func NonGenHttp() ServerTmplOption {
	return func(st *serverTmpl) {
		st.GenHttp = false
	}
}

func NewServerTmpl(appPkgPath, serverPath string, options ...ServerTmplOption) *serverTmpl {
	st := &serverTmpl{
		AppPkgPath: appPkgPath,
		ServerPath: serverPath,
		GenGrpc:    true,
		GenHttp:    true,
	}
	for _, option := range options {
		option(st)
	}
	return st
}

func (st serverTmpl) Generate() error {
	if st.GenGrpc {
		fmt.Println("generate server/grpc.go ...")
		grpcBuf := new(bytes.Buffer)
		tmpl, err := template.New("grpcServerTmpl").Parse(grpcServerTmpl)
		if err != nil {
			return err
		}
		if err := tmpl.Execute(grpcBuf, st); err != nil {
			return err
		}
		if err := os.WriteFile(path.Join(st.ServerPath, "grpc.go"), grpcBuf.Bytes(), 0o644); err != nil {
			return err
		}
	}
	if st.GenHttp {
		fmt.Println("generate server/http.go ...")
		httpBuf := new(bytes.Buffer)
		txTmpl, err := template.New("httpServerTmpl").Parse(httpServerTmpl)
		if err != nil {
			return err
		}
		if err := txTmpl.Execute(httpBuf, st); err != nil {
			return err
		}
		if err := os.WriteFile(path.Join(st.ServerPath, "http.go"), httpBuf.Bytes(), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func genServer(name, appPath string, isBff bool) error {
	// 4.1 mkdir server
	serverPath := path.Join(appPath, "server")
	if err := os.MkdirAll(serverPath, 0o700); err != nil {
		return err
	}

	var options []ServerTmplOption
	serverStr := "NewGRPCServer"
	if isBff {
		options = append(options, NonGenGrpc())
		serverStr = "NewHTTPServer"
	}

	// 4.2 gen server
	serverContent := `package server

import "github.com/google/wire"

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(%s)
`
	if err := os.WriteFile(path.Join(serverPath, "server.go"), []byte(fmt.Sprintf(serverContent, serverStr)), 0o644); err != nil {
		return err
	}
	// 4.3 gen grpc,http
	appPkgPath := path.Join(util.ModName(), "app", name)

	if err := NewServerTmpl(appPkgPath, serverPath, options...).Generate(); err != nil {
		return err
	}
	return nil
}
