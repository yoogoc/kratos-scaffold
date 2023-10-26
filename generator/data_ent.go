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
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
	"golang.org/x/tools/imports"
)

func (d *Data) EntFields() []*field.Field {
	if d.NeedAuditField {
		return util.FilterSlice(d.Fields, func(f *field.Field) bool {
			return f.Name != "updated_at" && f.Name != "created_at"
		})
	} else {
		return d.Fields
	}
}

//go:embed tmpl/data_ent_data.tmpl
var dataEntDataTmpl string

//go:embed tmpl/data_ent_schema.tmpl
var dataEntSchemaTmpl string

//go:embed tmpl/data_ent_transfer.tmpl
var dataEntTransferTmpl string

func (d *Data) GenerateEnt() error {
	// 1. gen ent schema and entity
	err := d.genEnt()
	if err != nil {
		return err
	}
	// 2. gen data transfer
	err = d.genTransferEnt()
	if err != nil {
		return err
	}
	// 3. gen data
	err = d.genDataEnt()
	if err != nil {
		return err
	}
	return nil
}

func (d *Data) EntPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	p := path.Join(wd, d.InternalPath(), "internal/data/ent")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := os.MkdirAll(p, 0o700); err != nil {
			panic(err)
		}
	}
	return p
}

func (d *Data) genEnt() error {
	fmt.Println("generating ent schema...")
	schemaBuf := new(bytes.Buffer)
	funcMap := template.FuncMap{
		"ToLower":  strings.ToLower,
		"ToPlural": util.Plural,
		"ToCamel":  strcase.ToCamel,
		"ToSnake":  strcase.ToSnake,
	}
	entSchemaTmpl, err := template.New("dataEntDataTmpl").Funcs(funcMap).Parse(dataEntSchemaTmpl)
	if err != nil {
		return err
	}
	err = entSchemaTmpl.Execute(schemaBuf, d)
	if err != nil {
		return err
	}
	p := path.Join(d.EntPath(), "schema", strings.ToLower(d.Name)+".go")
	if _, err := os.Stat(path.Join(d.EntPath(), "schema")); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Join(d.EntPath(), "schema"), 0o700); err != nil {
			return err
		}
	}
	entSchemaContent, err := imports.Process(p, schemaBuf.Bytes(), nil)
	if err != nil {
		return err
	}
	err = os.WriteFile(p, entSchemaContent, 0o644)
	if err != nil {
		return err
	}
	entGengoPath := path.Join(d.EntPath(), "generate.go")
	if _, err := os.Stat(entGengoPath); os.IsNotExist(err) {
		err := GenEntBase(d.EntPath())
		if err != nil {
			return err
		}
	}

	if err := util.Go("mod", "tidy"); err != nil {
		return err
	}
	fmt.Println("exec ent generate...")
	return util.Go("generate", d.EntPath())
}

func GenEntBase(entPath string) error {
	fmt.Println("generating ent generate...")
	content := `package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature privacy,entql,sql/modifier,sql/lock,sql/execquery,sql/upsert,namedges,sql/schemaconfig ./schema
`
	err := os.WriteFile(path.Join(entPath, "generate.go"), []byte(content), 0o644)
	if err != nil {
		return err
	}
	return nil
}

func (d *Data) genTransferEnt() error {
	fmt.Println("generating data transfer...")
	buf := new(bytes.Buffer)
	funcMap := template.FuncMap{
		"ToLower":      strings.ToLower,
		"ToPlural":     util.Plural,
		"ToCamel":      strcase.ToCamel,
		"ToLowerCamel": strcase.ToLowerCamel,
		"ToEntName":    field.EntName,
	}
	tmpl, err := template.New("dataEntTransferTmpl").Funcs(funcMap).Parse(dataEntTransferTmpl)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, d)
	if err != nil {
		return err
	}
	p := path.Join(d.OutPath(), strcase.ToSnake(d.Name)+"_transfer.go")
	content, err := imports.Process(p, buf.Bytes(), nil)
	if err != nil {
		return err
	}
	return os.WriteFile(p, content, 0o644)
}

func (d *Data) genDataEnt() error {
	fmt.Println("generating data...")
	buf := new(bytes.Buffer)
	funcMap := template.FuncMap{
		"ToLower":      strings.ToLower,
		"ToPlural":     util.Plural,
		"ToCamel":      strcase.ToCamel,
		"ToLowerCamel": strcase.ToLowerCamel,
		"ToEntName":    field.EntName,
		"last": func(x int, a []*field.Field) bool {
			return x == len(a)-1
		},
	}
	tmpl, err := template.New("dataEntDataTmpl").Funcs(funcMap).Parse(dataEntDataTmpl)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, d)
	if err != nil {
		return err
	}
	p := path.Join(d.OutPath(), strcase.ToSnake(d.Name)+".go")
	content, err := imports.Process(p, buf.Bytes(), nil)
	if err != nil {
		return err
	}
	return os.WriteFile(p, content, 0o644)
}
