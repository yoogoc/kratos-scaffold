package generator

import (
	"bytes"
	_ "embed"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/YoogoC/kratos-scaffold/pkg/cli"
	"github.com/YoogoC/kratos-scaffold/pkg/field"
	"github.com/YoogoC/kratos-scaffold/pkg/util"
)

type Proto struct {
	ApiDirName  string
	Namespace   string
	Name        string // is CamelName style
	Fields      field.Fields
	GenGrpc     bool
	GenHttp     bool
	primaryKey  string
	FieldStyle  string
	StrToPreMap map[string]field.PredicateType
}

func NewProto(setting *cli.EnvSettings) *Proto {
	return &Proto{
		primaryKey:  "id", // TODO
		FieldStyle:  setting.FieldStyle,
		StrToPreMap: field.StrToPreMap,
		ApiDirName:  setting.ApiDirName,
		Namespace:   setting.Namespace,
	}
}

func (p *Proto) CreateFields() []*field.Field {
	return p.Fields.CreateFields(p.primaryKey)
}

func (p *Proto) UpdateFields() []*field.Field {
	return p.Fields.UpdateFields(p.PrimaryField())
}

func (p *Proto) PrimaryField() *field.Field {
	idField, ok := util.FindSlice(p.Fields, func(f *field.Field) bool {
		return f.Name == p.primaryKey
	})
	if !ok {
		idField = &field.Field{
			Name:      p.primaryKey,
			FieldType: "int64", // TODO
			Predicates: []*field.Predicate{
				{
					Name:      p.primaryKey,
					Type:      field.PredicateTypeEq,
					FieldType: "int64",
				},
			},
		}
	}
	return idField
}

func (p *Proto) PageParamName() string {
	return field.StyleFieldMap[p.FieldStyle]("page")
}

func (p *Proto) PageSizeParamName() string {
	return field.StyleFieldMap[p.FieldStyle]("pageSize")
}

//go:embed tmpl/proto.tmpl
var protoTmpl string

func (p *Proto) Generate() error {
	buf := new(bytes.Buffer)

	funcMap := template.FuncMap{
		"ToLower":    strings.ToLower,
		"ToPlural":   util.Plural,
		"add":        func(a int, b int) int { return a + b },
		"fieldStyle": func(s string) string { return field.StyleFieldMap[p.FieldStyle](s) },
	}
	tmpl, err := template.New("protoTmpl").Funcs(funcMap).Parse(protoTmpl)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, p)
	if err != nil {
		return err
	}

	out := p.OutPath()
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if !filepath.IsAbs(out) {
		out = path.Join(wd, out)
	}
	n := strings.LastIndex(out, "/")
	op := out[:n]
	if _, err := os.Stat(op); os.IsNotExist(err) {
		if err := os.MkdirAll(op, 0o700); err != nil {
			return err
		}
	}
	return os.WriteFile(out, buf.Bytes(), 0o644)
}

func (p *Proto) OutPath() string {
	return path.Join(p.Path(), strings.ToLower(p.Name)+".proto")
}

func (p *Proto) Path() string {
	return path.Join(p.ApiDirName, p.Namespace, util.DefaultApiVersion)
}

func (p *Proto) GoPackage() string {
	s := strings.Split(p.Path(), "/")
	return util.ModName() + "/" + p.Path() + ";" + s[len(s)-1]
}

func (p *Proto) JavaPackage() string {
	return p.Package()
}

func (p *Proto) Package() string {
	return strings.ReplaceAll(p.Path(), "/", ".")
}
