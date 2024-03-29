package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/yoogoc/kratos-scaffold/pkg/cli"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"

	"github.com/iancoleman/strcase"
	"golang.org/x/tools/imports"
)

type Service struct {
	Base
	ApiPath string
	GenHttp bool
}

func NewService(setting *cli.EnvSettings) *Service {
	return &Service{
		Base: NewBase(setting, true),
	}
}

func (b *Service) CurrentPkgPath() string {
	return path.Join(util.ModName(), b.InternalPath(), "internal")
}

func (b *Service) FieldsExceptPrimary() []*field.Field {
	return util.FilterSlice(b.Fields, func(f *field.Field) bool {
		return f.Name != b.PrimaryKey
	})
}

func (b *Service) CreateFields() []*field.Field {
	return b.Fields.CreateFields(b.PrimaryKey)
}

func (b *Service) UpdateFields() []*field.Field {
	return b.Fields.UpdateFields(b.Fields.PrimaryField(b.PrimaryKey))
}

func (b *Service) PrimaryField() *field.Field {
	return b.Fields.PrimaryField(b.PrimaryKey)
}

//go:embed tmpl/service.tmpl
var serviceTmpl string

//go:embed tmpl/service_transfer.tmpl
var serviceTransferTmpl string

func (b *Service) Generate() error {
	if err := b.genTransfer(); err != nil {
		return err
	}
	if err := b.genService(); err != nil {
		return err
	}
	return nil
}

func (b *Service) genTransfer() error {
	fmt.Println("generating service transfer...")
	buf := new(bytes.Buffer)

	funcMap := template.FuncMap{
		"ToLower":    strings.ToLower,
		"ToPlural":   util.Plural,
		"ToCamel":    strcase.ToCamel,
		"ToSnake":    strcase.ToSnake,
		"ToLowCamel": strcase.ToLowerCamel,
	}
	tmpl, err := template.New("serviceTmpl").Funcs(funcMap).Parse(serviceTransferTmpl)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, b)
	if err != nil {
		return err
	}
	p := path.Join(b.OutPath(), strcase.ToSnake(b.Name)+"_transfer.go")
	content, err := imports.Process(p, buf.Bytes(), nil)
	if err != nil {
		return err
	}
	return os.WriteFile(p, content, 0o644)
}

func (b *Service) genService() error {
	fmt.Println("generating service...")
	buf := new(bytes.Buffer)

	funcMap := template.FuncMap{
		"ToLower":    strings.ToLower,
		"ToPlural":   util.Plural,
		"ToCamel":    strcase.ToCamel,
		"ToSnake":    strcase.ToSnake,
		"ToLowCamel": strcase.ToLowerCamel,
	}
	tmpl, err := template.New("serviceTmpl").Funcs(funcMap).Parse(serviceTmpl)
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

func (b *Service) OutPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(wd, b.InternalPath(), "internal/service") // TODO
}

func (b *Service) InternalPath() string {
	if b.Namespace != "" {
		return path.Join(b.AppDirName, b.Namespace) // TODO
	}
	return ""
}
