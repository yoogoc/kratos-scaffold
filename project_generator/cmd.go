package project_generator

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

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
