package project_generator

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/yoogoc/kratos-scaffold/pkg/util"
	"golang.org/x/tools/imports"
)

func genBffData(name, appPath string) error {
	dataPath := path.Join(appPath, "data")
	if err := os.MkdirAll(dataPath, 0o700); err != nil {
		return err
	}
	// 3.2 gen data/data.go data/tx.go
	appPkgPath := path.Join(util.ModName(), "app", name)
	if err := NewDataTmpl(name, appPkgPath, dataPath).Generate(dataGrpcTmpl); err != nil {
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

func (dt DataTmpl) Generate(text string) error {
	fmt.Println("generate data/data.go ...")
	dataBuf := new(bytes.Buffer)
	tmpl, err := template.New("dataTmpl").Parse(text)
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
	return nil
}
