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

type DataEnt struct {
	Base
	NeedAuditField bool
}

func NewDataEnt(setting *cli.EnvSettings) *DataEnt {
	return &DataEnt{
		Base: NewBase(setting, true),
	}
}

func (b *DataEnt) CreateFields() []*field.Field {
	return b.Fields.CreateFields(b.primaryKey)
}

func (b *DataEnt) EntFields() []*field.Field {
	if b.NeedAuditField {
		return util.FilterSlice(b.Fields, func(f *field.Field) bool {
			return f.Name != "updated_at" && f.Name != "created_at"
		})
	} else {
		return b.Fields
	}
}

//go:embed tmpl/data_ent_data.tmpl
var dataEntDataTmpl string

//go:embed tmpl/data_ent_schema.tmpl
var dataEntSchemaTmpl string

//go:embed tmpl/data_ent_transfer.tmpl
var dataEntTransferTmpl string

func (b *DataEnt) Generate() error {
	// 1. gen ent schema and entity
	err := b.genEnt()
	if err != nil {
		return err
	}
	// 2. gen data transfer
	err = b.genTransfer()
	if err != nil {
		return err
	}
	// 3. gen data
	err = b.genData()
	if err != nil {
		return err
	}
	return nil
}

func (b *DataEnt) EntPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	p := path.Join(wd, b.InternalPath(), "internal/data/ent")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := os.MkdirAll(p, 0o700); err != nil {
			panic(err)
		}
	}
	return p
}

func (b *DataEnt) OutPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(wd, b.InternalPath(), "internal/data")
}

func (b *DataEnt) CurrentPkgPath() string {
	return path.Join(util.ModName(), b.InternalPath(), "internal")
}

func (b *DataEnt) genEnt() error {
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
	err = entSchemaTmpl.Execute(schemaBuf, b)
	if err != nil {
		return err
	}
	p := path.Join(b.EntPath(), "schema", strings.ToLower(b.Name)+".go")
	if _, err := os.Stat(path.Join(b.EntPath(), "schema")); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Join(b.EntPath(), "schema"), 0o700); err != nil {
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
	entGengoPath := path.Join(b.EntPath(), "generate.go")
	if _, err := os.Stat(entGengoPath); os.IsNotExist(err) {
		err := GenEntBase(b.EntPath())
		if err != nil {
			return err
		}
	}

	if err := util.Go("mod", "tidy"); err != nil {
		return err
	}
	fmt.Println("exec ent generate...")
	return util.Go("generate", b.EntPath())
}

func GenEntBase(entPath string) error {
	if err := os.MkdirAll(path.Join(entPath, "external"), 0o700); err != nil {
		return err
	}
	fmt.Println("generating ent external...")
	externalContent := `{{ define "external" }}
package ent

import (
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

func (c *Client) DB() *sql.DB {
	switch d := c.driver.(type) {
	case *entsql.Driver:
		return d.DB()
	case *dialect.DebugDriver:
		return d.Driver.(*entsql.Driver).DB()
	default:
		panic("unknown driver")
	}
}

func (c *Client) Driver() dialect.Driver {
	return c.driver
}
{{ end }}
`
	err := os.WriteFile(path.Join(entPath, "external", "sql.tmpl"), []byte(externalContent), 0o644)
	if err != nil {
		return err
	}
	fmt.Println("generating ent generate...")
	content := `package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --template ./external --feature privacy,sql/modifier,sql/lock ./schema
`
	err = os.WriteFile(path.Join(entPath, "generate.go"), []byte(content), 0o644)
	if err != nil {
		return err
	}
	return nil
}

func (b *DataEnt) genTransfer() error {
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

func (b *DataEnt) genData() error {
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

func (b *DataEnt) InternalPath() string {
	if b.Namespace != "" {
		return path.Join(b.AppDirName, b.Namespace)
	}
	return ""
}
