package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/yoogoc/kratos-scaffold/generator/wire"
	"github.com/yoogoc/kratos-scaffold/pkg/cli"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
	"golang.org/x/tools/imports"
)

type Biz struct {
	Base
	Http bool
}

func NewBiz(setting *cli.EnvSettings) *Biz {
	return &Biz{
		Base: NewBase(setting, true),
	}
}

func (b *Biz) PrimaryField() *field.Field {
	return b.Fields.PrimaryField(b.PrimaryKey)
}

func (b *Biz) PatchFields() []*field.Field {
	return b.Fields.PatchFields(b.Fields.PrimaryField(b.PrimaryKey))
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
	if err := os.WriteFile(p, content, 0o644); err != nil {
		return err
	}

	providerSetPath := path.Join(b.OutPath(), "biz.go")
	if err := wire.AddToProviderSet(providerSetPath, "New"+b.Name+"Usecase"); err != nil {
		fmt.Printf("failed to add to ProviderSet, please add manually: %v\n", err)
	}
	return nil
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
