package generator

import (
	"bytes"
	_ "embed"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/YoogoC/kratos-scaffold/pkg/util"

	"github.com/iancoleman/strcase"
	"golang.org/x/tools/imports"
)

type Service struct {
	Name       string
	Namespace  string
	AppDirName string // TODO
	Fields     []Field
	ApiPath    string
}

func NewService(name string, ns string, fields []Field) Service {
	adn := ""
	apiModelName := name
	if ns != "" {
		adn = "app/" + ns // TODO
		apiModelName = ns
	}

	return Service{
		Name:       plural.Singular(strings.ToUpper(name[0:1]) + name[1:]),
		Namespace:  ns,
		AppDirName: adn,
		Fields:     fields,
		ApiPath:    path.Join(ModName(), "api", apiModelName, "v1"), // TODO
	}
}

func (b Service) ParamFields() []Field {
	fs := make([]Field, 0, len(b.Fields))
	for _, field := range b.Fields {
		for _, predicate := range field.Predicates {
			fs = append(fs, Field{
				Name:      field.Name + predicate.Type.String(),
				FieldType: field.FieldType,
			})
		}
	}
	return fs
}

func (b Service) CurrentPkgPath() string {
	return path.Join(ModName(), b.AppDirName, "internal")
}

func (b Service) FieldsExceptPrimary() []Field {
	return util.FilterSlice(b.Fields, func(f Field) bool {
		return f.Name != "id"
	})
}

//go:embed tmpl/service.tmpl
var serviceTmpl string

//go:embed tmpl/service_transfer.tmpl
var serviceTransferTmpl string

func (b Service) Generate() error {
	if err := b.genTransfer(); err != nil {
		return err
	}
	if err := b.genService(); err != nil {
		return err
	}
	return nil
}

func (b Service) genTransfer() error {
	buf := new(bytes.Buffer)

	funcMap := template.FuncMap{
		"ToLower":    strings.ToLower,
		"ToPlural":   plural.Plural,
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

func (b Service) genService() error {
	buf := new(bytes.Buffer)

	funcMap := template.FuncMap{
		"ToLower":    strings.ToLower,
		"ToPlural":   plural.Plural,
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

func (b Service) OutPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(wd, b.AppDirName, "internal/service") // TODO
}
