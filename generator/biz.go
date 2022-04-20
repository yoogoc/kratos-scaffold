package generator

import (
	"bytes"
	_ "embed"
	"os"
	"path"
	"strings"
	"text/template"
	
	"github.com/iancoleman/strcase"
	"golang.org/x/tools/imports"
)

type Biz struct {
	Name        string
	Namespace   string
	AppDirName  string // TODO
	Fields      []Field
	StrToPreMap map[string]PredicateType
}

func NewBiz(name string, ns string, fields []Field) Biz {
	adn := ""
	if ns != "" {
		adn = "app/" + ns // TODO
	}
	return Biz{
		Name:        plural.Singular(strings.ToUpper(name[0:1]) + name[1:]),
		Namespace:   ns,
		AppDirName:  adn,
		Fields:      fields,
		StrToPreMap: strToPreMap,
	}
}

func (b Biz) ParamFields() []Predicate {
	fs := make([]Predicate, 0, len(b.Fields))
	for _, field := range b.Fields {
		for _, predicate := range field.Predicates {
			fs = append(fs, Predicate{
				Name:      field.Name + predicate.Type.String(),
				FieldType: field.FieldType,
				Type:      predicate.Type,
			})
		}
	}
	return fs
}

//go:embed tmpl/biz.tmpl
var bizTmpl string

func (b Biz) Generate() error {
	buf := new(bytes.Buffer)
	
	funcMap := template.FuncMap{
		"ToLower":  strings.ToLower,
		"ToPlural": plural.Plural,
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

func (b Biz) OutPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(wd, b.AppDirName, "internal/biz") // TODO
}
