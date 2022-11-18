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

//go:embed tmpl/data_proto_data.tmpl
var dataProtoDataTmpl string

//go:embed tmpl/data_proto_transfer.tmpl
var dataProtoTransferTmpl string

func (d *Data) GenerateProto() error {
	// 2. gen data transfer

	if err := d.genTransferProto(); err != nil {
		return err
	}
	// 3. gen data
	if err := d.genDataProto(); err != nil {
		return err
	}
	return nil
}

func (d *Data) genTransferProto() error {
	fmt.Println("generating data transfer...")
	buf := new(bytes.Buffer)
	funcMap := template.FuncMap{
		"ToLower":      strings.ToLower,
		"ToPlural":     util.Plural,
		"ToCamel":      strcase.ToCamel,
		"ToLowerCamel": strcase.ToLowerCamel,
		"ToEntName":    field.EntName,
	}
	tmpl, err := template.New("dataProtoTransferTmpl").Funcs(funcMap).Parse(dataProtoTransferTmpl)
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

func (d *Data) genDataProto() error {
	fmt.Println("generating data...")
	buf := new(bytes.Buffer)
	funcMap := template.FuncMap{
		"ToLower":       strings.ToLower,
		"ToPlural":      util.Plural,
		"ToCamel":       strcase.ToCamel,
		"ToCamelPlural": func(s string) string { return strcase.ToCamel(util.Plural(s)) },
		"ToLowerCamel":  strcase.ToLowerCamel,
		"last": func(x int, a []*field.Field) bool {
			return x == len(a)-1
		},
	}
	tmpl, err := template.New("dataProtoDataTmpl").Funcs(funcMap).Parse(dataProtoDataTmpl)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, d)
	if err != nil {
		return err
	}
	p := path.Join(d.OutPath(), strcase.ToSnake(d.Name)+".go")
	// content := buf.Bytes()
	content, err := imports.Process(p, buf.Bytes(), nil)
	if err != nil {
		return err
	}
	return os.WriteFile(p, content, 0o644)
}
