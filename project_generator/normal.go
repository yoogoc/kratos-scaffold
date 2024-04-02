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
	bizContent := `package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewTxUsecase)
`
	if err := os.WriteFile(path.Join(bizPath, "biz.go"), []byte(bizContent), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(bizPath, "tx.go"), []byte(bizTxGo), 0o644); err != nil {
		return err
	}
	// 2. mkdir service. gen service/service.go
	servicePath := path.Join(appPath, "service")
	if err := os.MkdirAll(servicePath, 0o700); err != nil {
		return err
	}
	serviceContent := `package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet()
`
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
	serverContent := `package server

import "github.com/google/wire"

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)
`
	if err := os.WriteFile(path.Join(serverPath, "server.go"), []byte(serverContent), 0o644); err != nil {
		return err
	}
	// 4.3 gen grpc,http
	if err := NewServerTmpl(appPkgPath, serverPath).Generate(); err != nil {
		return err
	}

	// 5 gen conf
	confPath := path.Join(appPath, "conf")
	if err := os.MkdirAll(confPath, 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(confPath, "conf.proto"), []byte(confProto), 0o644); err != nil {
		return err
	}

	// if err := util.Exec("make", "config-"+name); err != nil {
	// 	return err
	// }

	// 6 gen log
	logPath := path.Join(appPath, "log")
	if err := os.MkdirAll(logPath, 0o700); err != nil {
		return err
	}
	if err := NewLogTmpl(appPkgPath, logPath).Generate(); err != nil {
		return err
	}

	return nil
}
