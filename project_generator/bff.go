package project_generator

import (
	"os"
	"path"
	"strings"

	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

func genBffInternal(name string, appPath string, isSubMono bool) error {
	// biz,data,service,server,conf
	// orm := "grpc"
	// 1. mkdir biz. gen biz/biz.go
	if err := genBffBiz(appPath); err != nil {
		return err
	}

	// 2. mkdir service. gen service/service.go
	if err := genService(appPath); err != nil {
		return err
	}

	// 3. gen data
	if err := genBffData(name, appPath); err != nil {
		return err
	}

	// 4 mkdir server. gen server
	if err := genServer(name, appPath, true); err != nil {
		return err
	}

	// 5 gen conf
	if err := genConf(name, appPath, confBffProto); err != nil {
		return err
	}

	// 6 gen log
	logPath := path.Join(appPath, "log")
	if err := os.MkdirAll(logPath, 0o700); err != nil {
		return err
	}
	appPkgPath := name
	if isSubMono {
		appPkgPath = path.Join(util.ModName(), "app", name)
	}
	if err := NewLogTmpl(appPkgPath, logPath).Generate(); err != nil {
		return err
	}

	// 7 gen otel
	otelPath := path.Join(appPath, "otel")
	if err := os.MkdirAll(otelPath, 0o700); err != nil {
		return err
	}
	serviceName := strings.ReplaceAll(appPkgPath, "/", ".")
	if err := NewOtelTmpl(appPkgPath, serviceName, otelPath).Generate(); err != nil {
		return err
	}

	// 8 gen middleware
	middlewarePath := path.Join(appPath, "middleware")
	if err := os.MkdirAll(middlewarePath, 0o700); err != nil {
		return err
	}
	middlewareContent := `package middleware

import "github.com/google/wire"

var ProviderSet = wire.NewSet()
`
	if err := os.WriteFile(path.Join(middlewarePath, "middleware.go"), []byte(middlewareContent), 0o644); err != nil {
		return err
	}

	// 9 gen demo greeter
	if err := genDemo(name, appPath, isSubMono); err != nil {
		return err
	}

	return nil
}
