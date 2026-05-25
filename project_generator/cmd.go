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

func genCmd(name string, appPath string, isSubMono, isBff bool) error {
	// mkdir cmd/server. gen main.go, wire.go, wire_gen.go
	mainPath := path.Join(appPath, "cmd")
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
		IsBff:       isBff,
		IsSubMono:   isSubMono,
		Name:        name,
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
	IsBff       bool
	IsSubMono   bool
	Name        string
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
	if err := util.WhiteGo(path.Join(ct.OutPath, "main.go"), mainBuf.Bytes()); err != nil {
		return err
	}

	fmt.Println("generate cmd/server/server.go ...")
	serverBuf := new(bytes.Buffer)
	serverTmpl, err := template.New("serverTmpl").Parse(cmdServerTmpl)
	if err != nil {
		return err
	}
	if err := serverTmpl.Execute(serverBuf, ct); err != nil {
		return err
	}
	if err := util.WhiteGo(path.Join(ct.OutPath, "server.go"), serverBuf.Bytes()); err != nil {
		return err
	}

	if !ct.IsBff {
		fmt.Println("generate cmd/server/migration.go ...")
		migrationBuf := new(bytes.Buffer)
		migrationTmpl, err := template.New("migrationTmpl").Parse(cmdMigrationTmpl)
		if err != nil {
			return err
		}
		if err := migrationTmpl.Execute(migrationBuf, ct); err != nil {
			return err
		}
		if err := util.WhiteGo(path.Join(ct.OutPath, "migration.go"), migrationBuf.Bytes()); err != nil {
			return err
		}
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
	if err := util.WhiteGo(path.Join(ct.OutPath, "wire.go"), wireBuf.Bytes()); err != nil {
		return err
	}

	if _, err := exec.LookPath("wire"); err != nil {
		fmt.Println("wire not found, installing...")
		if err := util.Go("install", "github.com/google/wire/cmd/wire@latest"); err != nil {
			fmt.Println("wire install failed, please install manually: go install github.com/google/wire/cmd/wire@latest")
			return nil
		}
	}
	fmt.Println("running wire generate ...")
	if ct.IsSubMono {
		if err := util.Exec("make", "wire-"+ct.Name); err != nil {
			fmt.Printf("wire generate failed, please run manually: make wire-%s\n", ct.Name)
		}
	} else {
		if err := util.Exec("wire", path.Join(".", ct.OutPath)); err != nil {
			fmt.Printf("wire generate failed, please run manually: wire ./%s\n", ct.OutPath)
		}
	}

	return nil
}
