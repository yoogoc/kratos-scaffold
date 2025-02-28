package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

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
	return util.Plural(strcase.ToSnake(d.Name))
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

	fileSuffix := fmt.Sprintf("_create_%s.sql", util.Plural(strcase.ToSnake(d.Name)))

	entries, err := os.ReadDir(d.MigrationPath())
	if err != nil {
		return err
	}
	var needRemoveFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), fileSuffix) {
			needRemoveFiles = append(needRemoveFiles, path.Join(d.MigrationPath(), entry.Name()))
		}
	}
	for _, file := range needRemoveFiles {
		_ = os.Remove(file)
	}

	p := path.Join(d.MigrationPath(), fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), fileSuffix))

	if err := os.WriteFile(p, schemaBuf.Bytes(), 0o644); err != nil {
		return err
	}
	return nil
}
