package project_generator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

type DemoTmpl struct {
	AppPkgPath   string
	ApiPkgPath   string
	ProtoPackage string
	GoPackage    string
}

func genDemo(name, appPath string, isSubMono bool) error {
	wd, _ := os.Getwd()

	var apiProtoDir string
	var apiPkgPath string
	var protoPackage string

	if isSubMono {
		apiProtoDir = path.Join(wd, "api", name, "v1")
		apiPkgPath = path.Join(util.ModName(), "api", name, "v1")
		protoPackage = strings.ReplaceAll(strings.ReplaceAll(path.Join("api", name, "v1"), "/", "."), "-", "_")
	} else {
		apiProtoDir = path.Join(appPath, "..", "api", "v1")
		apiPkgPath = path.Join(name, "api", "v1")
		protoPackage = strings.ReplaceAll(strings.ReplaceAll(path.Join("api", "v1"), "/", "."), "-", "_")
	}
	goPackage := apiPkgPath + ";v1"

	appPkgPath := name
	if isSubMono {
		appPkgPath = path.Join(util.ModName(), "app", name)
	}

	dt := DemoTmpl{
		AppPkgPath:   appPkgPath,
		ApiPkgPath:   apiPkgPath,
		ProtoPackage: protoPackage,
		GoPackage:    goPackage,
	}

	// 1. mkdir api proto dir
	if err := os.MkdirAll(apiProtoDir, 0o700); err != nil {
		return err
	}

	// 2. render greeter.proto
	fmt.Println("generating demo greeter.proto ...")
	protoBuf := new(bytes.Buffer)
	protoTmpl, err := template.New("demoProtoTmpl").Parse(demoGreeterProtoTmpl)
	if err != nil {
		return err
	}
	if err := protoTmpl.Execute(protoBuf, dt); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(apiProtoDir, "greeter.proto"), protoBuf.Bytes(), 0o644); err != nil {
		return err
	}

	// 3. compile proto
	fmt.Println("compiling demo greeter.proto ...")
	if isSubMono {
		if err := util.Exec("make", "proto-"+name); err != nil {
			fmt.Printf("proto compile failed, please run manually: make proto-%s\n", name)
		}
	} else {
		protoFile := path.Join("api", "v1", "greeter.proto")
		cmd := exec.Command("protoc",
			"--proto_path=.",
			"--proto_path=./third_party",
			"--go_out=paths=source_relative:.",
			"--go-grpc_out=paths=source_relative:.",
			"--go-http_out=paths=source_relative:.",
			protoFile,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("proto compile failed, please run manually: protoc ... %s\n", protoFile)
		}
	}

	// 4. write biz/greeter.go (static content)
	fmt.Println("generating demo biz/greeter.go ...")
	if err := os.WriteFile(path.Join(appPath, "biz", "greeter.go"), []byte(demoGreeterBizGo), 0o644); err != nil {
		return err
	}

	// 5. render service/greeter.go
	fmt.Println("generating demo service/greeter.go ...")
	serviceBuf := new(bytes.Buffer)
	serviceTmpl, err := template.New("demoServiceTmpl").Parse(demoGreeterServiceTmpl)
	if err != nil {
		return err
	}
	if err := serviceTmpl.Execute(serviceBuf, dt); err != nil {
		return err
	}
	if err := util.WhiteGo(path.Join(appPath, "service", "greeter.go"), serviceBuf.Bytes()); err != nil {
		return err
	}

	// 6. render data/greeter.go
	fmt.Println("generating demo data/greeter.go ...")
	dataBuf := new(bytes.Buffer)
	dataTmpl, err := template.New("demoDataTmpl").Parse(demoGreeterDataTmpl)
	if err != nil {
		return err
	}
	if err := dataTmpl.Execute(dataBuf, dt); err != nil {
		return err
	}
	if err := util.WhiteGo(path.Join(appPath, "data", "greeter.go"), dataBuf.Bytes()); err != nil {
		return err
	}

	return nil
}
