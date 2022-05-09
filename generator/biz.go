package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/YoogoC/kratos-scaffold/pkg/cli"
	"github.com/YoogoC/kratos-scaffold/pkg/util"
	"github.com/iancoleman/strcase"
	"golang.org/x/tools/imports"
)

type Biz struct {
	Base
}

func NewBiz(setting *cli.EnvSettings) *Biz {
	return &Biz{
		Base: NewBase(setting, true),
	}
}

//go:embed tmpl/biz.tmpl
var bizTmpl string

func (b *Biz) Generate() error {
	fmt.Println("generating biz...")
	buf := new(bytes.Buffer)

	funcMap := template.FuncMap{
		"ToLower":  strings.ToLower,
		"ToPlural": util.Plural,
		"ToCamel":  strcase.ToCamel,
	}
	tmpl, err := template.New("bizTmpl").Funcs(funcMap).Parse(bizTmpl)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, b)
	if err != nil {
		return err
	}
	p := path.Join(b.OutPath(), strcase.ToSnake(b.Name)+".go")
	content, err := imports.Process(p, buf.Bytes(), nil)
	if err != nil {
		return err
	}
	return os.WriteFile(p, content, 0o644)
}

func (b *Biz) OutPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if b.Namespace != "" {
		return path.Join(wd, b.AppDirName, b.Namespace, "internal/biz") // TODO
	} else {
		return path.Join(wd, b.Namespace, "internal/biz") // TODO
	}

}
