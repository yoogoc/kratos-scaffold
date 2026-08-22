package project_generator

import (
	"os"
	"path"

	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

func genInternal(name string, appPath string, isSubMono bool) error {
	// biz,data,service,server,conf
	// orm := "ent"
	// 1. mkdir biz. gen biz/biz.go, biz/tx.go
	bizPath := path.Join(appPath, "biz")
	if err := os.MkdirAll(bizPath, 0o700); err != nil {
		return err
	}
	bizContent := providerDeclaration("biz", "NewTxUsecase", "NewGreeterUsecase")
	if err := os.WriteFile(path.Join(bizPath, "biz.go"), []byte(bizContent), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(bizPath, "tx.go"), []byte(bizTxGo), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(bizPath, "pagination.go"), []byte(kitPaginationAliases), 0o644); err != nil {
		return err
	}
	// 2. mkdir service. gen service/service.go
	servicePath := path.Join(appPath, "service")
	if err := os.MkdirAll(servicePath, 0o700); err != nil {
		return err
	}
	serviceContent := providerDeclaration("service", "NewGreeterService")
	if err := os.WriteFile(path.Join(servicePath, "service.go"), []byte(serviceContent), 0o644); err != nil {
		return err
	}
	// 3. mkdir data, init ent, gen data/data.go data/tx.go
	// 3.1 mkdir data
	dataPath := path.Join(appPath, "data")
	if err := os.MkdirAll(dataPath, 0o700); err != nil {
		return err
	}
	// 3.2 init ent
	if err := initEnt(dataPath); err != nil {
		return err
	}
	// 3.3 gen data/data.go data/tx.go
	appPkgPath := name
	if isSubMono {
		appPkgPath = path.Join(util.ModName(), "app", name)
	}
	if err := NewDataTmpl(name, appPkgPath, dataPath).Generate(dataEntTmpl); err != nil {
		return err
	}
	// 4 mkdir server. gen server
	// 4.1 mkdir server
	serverPath := path.Join(appPath, "server")
	if err := os.MkdirAll(serverPath, 0o700); err != nil {
		return err
	}
	// 4.2 gen server
	serverContent := providerDeclaration("server", "NewGRPCServer", "NewHTTPServer")
	if err := os.WriteFile(path.Join(serverPath, "server.go"), []byte(serverContent), 0o644); err != nil {
		return err
	}
	// 4.3 gen grpc,http
	apiPkgPath := appPkgPath
	if isSubMono {
		apiPkgPath = path.Join(util.ModName(), "api", name, "v1")
	} else {
		apiPkgPath = path.Join(name, "api", "v1")
	}
	if err := NewServerTmpl(appPkgPath, apiPkgPath, serverPath).Generate(); err != nil {
		return err
	}

	// 5 gen conf + buf generate
	if isSubMono {
		if err := genConf(name, appPath, confProto); err != nil {
			return err
		}
	} else {
		if err := genConfSingle(appPath, confProto); err != nil {
			return err
		}
	}

	// 6 gen middleware
	middlewarePath := path.Join(appPath, "middleware")
	if err := os.MkdirAll(middlewarePath, 0o700); err != nil {
		return err
	}
	middlewareContent := providerDeclaration("middleware")
	if err := os.WriteFile(path.Join(middlewarePath, "middleware.go"), []byte(middlewareContent), 0o644); err != nil {
		return err
	}

	// 7 gen demo greeter
	if err := genDemo(name, appPath, isSubMono, false); err != nil {
		return err
	}
	return genDoAdapters(appPkgPath, appPath, false)
}
