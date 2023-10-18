package project_generator

import (
	"os"
	"path"

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
	return nil
}
