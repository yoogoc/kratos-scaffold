package project_generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/yoogoc/kratos-scaffold/generator"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
	"golang.org/x/tools/imports"
)

func GenMono(name string) error {
	// 1. find project name
	ss := strings.Split(name, "/")
	projectName := ss[len(ss)-1]
	// 2. mkdir <project name>
	wd, _ := os.Getwd()
	projectPath := path.Join(wd, projectName)
	if err := os.MkdirAll(projectPath, 0o700); err != nil {
		return err
	}
	// 3. Chdir <project path>
	err := os.Chdir(projectPath)
	if err != nil {
		return err
	}
	// 4. init go.mod, readme
	if err = genGoMod(name, projectPath); err != nil {
		return err
	}

	err = util.Go(
		"get",
		"github.com/go-kratos/kratos/v2",
		"google.golang.org/grpc",
		"google.golang.org/protobuf",
		"github.com/google/wire",
		"github.com/pkg/errors",
		"github.com/gorilla/handlers",
	)
	if err != nil {
		return err
	}
	// fd := exec.Command("go", "mod", "tidy")
	// fd.Stdout = os.Stdout
	// fd.Stderr = os.Stderr
	// fd.Dir = "."
	// if err = fd.Run(); err != nil {
	// 	return err
	// }
	if err = os.WriteFile(path.Join(projectPath, "README.md"), []byte(name), 0o644); err != nil {
		return err
	}
	// 5. mkdir api, app, pkg
	fmt.Println("generating api/ ...")
	if err := util.GenNullPath(path.Join(projectPath, "api")); err != nil {
		return err
	}
	fmt.Println("generating app/ ...")
	if err := util.GenNullPath(path.Join(projectPath, "app")); err != nil {
		return err
	}
	fmt.Println("generating pkg/ ...")
	if err := util.GenNullPath(path.Join(projectPath, "pkg")); err != nil {
		return err
	}
	// 6. cp common proto
	if err := cpProto(projectPath); err != nil {
		return err
	}
	return nil
}

func Gen(name string) error {
	// 1. 判断该生成单体服务还是mono下子服务,目前通过当前目录是否存在go.mod判定
	isMono := true
	wd, _ := os.Getwd()
	if _, err := os.Stat(path.Join(wd, "go.mod")); err != nil {
		if os.IsNotExist(err) {
			isMono = false
		} else {
			return err
		}
	}
	if isMono {
		// name 只能是 xx
		return GenSubMono(name)
	} else {
		// name 可能是 xx, aa.com/xx/yy
		return GenSingle(name)
	}

}

func GenSubMono(name string) error {
	appDirName := "app"
	// 1. gen sub app path
	wd, _ := os.Getwd()
	subAppPath := path.Join(wd, appDirName, name)
	if err := os.MkdirAll(subAppPath, 0o700); err != nil {
		return err
	}
	// 2. gen internal: biz,data,service,server,conf
	if err := genInternal(name, path.Join(subAppPath, "internal"), true); err != nil {
		return err
	}
	// 3. gen cmd
	if err := genCmd(name, subAppPath, true); err != nil {
		return err
	}
	// 4. init configs/conf.yaml
	if err := genConfigs(name, subAppPath); err != nil {
		return err
	}
	return nil
}

func GenSingle(name string) error {
	// 1. mkdir app path
	wd, _ := os.Getwd()
	appPath := path.Join(wd, name)
	if err := os.MkdirAll(appPath, 0o700); err != nil {
		return err
	}
	if err := os.Chdir(appPath); err != nil {
		return err
	}
	// 2. cp common proto
	if err := cpProto(appPath); err != nil {
		return err
	}
	// 3. gen go.mod
	if err := genGoMod(name, appPath); err != nil {
		return err
	}
	// 4. gen internal: biz,data,service,server,conf
	if err := genInternal(name, path.Join(appPath, "internal"), false); err != nil {
		return err
	}
	// 5. gen cmd
	if err := genCmd(name, appPath, false); err != nil {
		return err
	}
	// 6. init configs/conf.yaml
	if err := genConfigs(name, appPath); err != nil {
		return err
	}
	return nil
}

func genGoMod(name, projectPath string) error {
	goModContent := fmt.Sprintf("module %s\ngo 1.18", name)
	if err := os.WriteFile(path.Join(projectPath, "go.mod"), []byte(goModContent), 0o644); err != nil {
		return err
	}
	return nil
}

func genConfigs(name string, appPath string) error {
	configsPath := path.Join(appPath, "configs")
	if err := os.MkdirAll(configsPath, 0o700); err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(configsPath, "config.yaml"), []byte(configContent), 0o644); err != nil {
		return err
	}
	return nil
}

func genCmd(name string, appPath string, isSubMono bool) error {
	// mkdir cmd/server. gen main.go, wire.go, wire_gen.go
	mainPath := path.Join(appPath, "cmd/server")
	if err := os.MkdirAll(mainPath, 0o700); err != nil {
		return err
	}
	appPkgPath := name
	serviceName := strings.ReplaceAll(name, "/", ".")
	if isSubMono {
		appPkgPath = path.Join(util.ModName(), "app", name)
		serviceName = strings.ReplaceAll(path.Join(util.ModName(), name), "/", ".")
	}
	cmdTmpl := CmdTmpl{
		AppPkgPath:  appPkgPath,
		ServiceName: serviceName,
		OutPath:     mainPath,
	}
	if err := cmdTmpl.Generate(); err != nil {
		return err
	}
	return nil
}

type CmdTmpl struct {
	AppPkgPath  string
	ServiceName string
	OutPath     string
}

func (ct CmdTmpl) Generate() error {
	fmt.Println("generate cmd/server/main.go ...")
	mainBuf := new(bytes.Buffer)
	tmpl, err := template.New("mainTmpl").Parse(cmdMainTmpl)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(mainBuf, ct); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(ct.OutPath, "main.go"), mainBuf.Bytes(), 0o644); err != nil {
		return err
	}

	fmt.Println("generate cmd/server/wire.go ...")
	wireBuf := new(bytes.Buffer)
	wireTmpl, err := template.New("wireTmpl").Parse(cmdWireTmpl)
	if err != nil {
		return err
	}
	if err := wireTmpl.Execute(wireBuf, ct); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(ct.OutPath, "wire.go"), wireBuf.Bytes(), 0o644); err != nil {
		return err
	}

	// if err := util.Go("get", "github.com/google/wire"); err != nil {
	// 	return err
	// }
	// if err := util.Go("mod", "tidy"); err != nil {
	// 	return err
	// }
	// fd := exec.Command("wire", ct.OutPath)
	// fd.Stdout = os.Stdout
	// fd.Stderr = os.Stderr
	// fd.Dir = "."
	// if err := fd.Run(); err != nil {
	// 	return err
	// }
	fmt.Println("skip wire generate")

	return nil
}

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
	if err := NewDataTmpl(name, appPkgPath, dataPath).Generate(); err != nil {
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
	confProtoPath := path.Join("internal/conf", "conf.proto")
	if isSubMono {
		confProtoPath = path.Join("app", name, confProtoPath)
	}
	fd := exec.Command("protoc",
		"--proto_path=.",
		"--proto_path=./third_party",
		"--go_out=paths=source_relative:.",
		path.Join(".", confProtoPath),
	)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if err := fd.Run(); err != nil {
		return err
	}
	return nil
}

type ServerTmpl struct {
	AppPkgPath string
	ServerPath string
}

func NewServerTmpl(appPkgPath, serverPath string) ServerTmpl {
	return ServerTmpl{
		AppPkgPath: appPkgPath,
		ServerPath: serverPath,
	}
}

func (st ServerTmpl) Generate() error {
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
	return nil
}

type DataTmpl struct {
	AppPkgPath   string
	DBDriverType string
	LoggerName   string
	DataPath     string
}

func NewDataTmpl(name, appPkgPath, dataPath string) DataTmpl {
	return DataTmpl{
		AppPkgPath:   appPkgPath,
		DBDriverType: "mysql",
		LoggerName:   name,
		DataPath:     dataPath,
	}
}

func (dt DataTmpl) Generate() error {
	fmt.Println("generate data/data.go ...")
	dataBuf := new(bytes.Buffer)
	tmpl, err := template.New("dataTmpl").Parse(dataTmpl)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(dataBuf, dt); err != nil {
		return err
	}
	dp := path.Join(dt.DataPath, "data.go")
	content, err := imports.Process(dp, dataBuf.Bytes(), nil)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dp, content, 0o644); err != nil {
		return err
	}
	fmt.Println("generate data/tx.go ...")
	return nil
}

func initEnt(appDataPath string) error {
	entPath := path.Join(appDataPath, "ent")
	entSchemaPath := path.Join(entPath, "schema")
	if err := os.MkdirAll(entSchemaPath, 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(entSchemaPath, "needremove.go"), []byte(entSchemaNeedRemove), 0o644); err != nil {
		return err
	}
	if err := generator.GenEntBase(entPath); err != nil {
		return err
	}
	if err := util.Go("generate", entPath); err != nil {
		return err
	}
	fmt.Println("ent至少有一个schema才会生成客户端代码,项目生成后,请自行删除ent/下needremove开头的文件和文件夹")
	return nil
}
