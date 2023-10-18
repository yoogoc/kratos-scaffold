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
	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

//go:embed tmpl/migration.sql.tmpl
var dataMigrationTmpl string

func (d *Data) MigrationPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	p := path.Join(wd, d.InternalPath(), "db/migration")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := os.MkdirAll(p, 0o700); err != nil {
			panic(err)
		}
	}
	return p
}

func (d *Data) TableName() string {
	return util.Plural(strings.ToLower(d.Name))
}
func (d *Data) GenerateMigration() error {
	fmt.Println("generating migration...")
	schemaBuf := new(bytes.Buffer)
	funcMap := template.FuncMap{
		"ToLower":  strings.ToLower,
		"ToPlural": util.Plural,
		"ToCamel":  strcase.ToCamel,
		"ToSnake":  strcase.ToSnake,
	}
	entSchemaTmpl, err := template.New("dataMigrationTmpl").Funcs(funcMap).Parse(dataMigrationTmpl)
	if err != nil {
		return err
	}
	err = entSchemaTmpl.Execute(schemaBuf, d)
	if err != nil {
		return err
	}
	p := path.Join(d.MigrationPath(), fmt.Sprintf("create_%s.sql", util.Plural(strings.ToLower(d.Name))))

	if err := os.WriteFile(p, schemaBuf.Bytes(), 0o644); err != nil {
		return err
	}
	return nil
}
