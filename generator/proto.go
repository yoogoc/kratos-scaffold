package generator

import (
	"bytes"
	_ "embed"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/YoogoC/kratos-scaffold/pkg/field"
	"github.com/YoogoC/kratos-scaffold/pkg/util"

	"github.com/gertd/go-pluralize"
	"golang.org/x/mod/modfile"
)

type Proto struct {
	Name        string
	OutPath     string
	Package     string
	GoPackage   string
	JavaPackage string
	Fields      []field.Field
	GenGrpc     bool
	GenHttp     bool
	primaryKey  string
	fieldStyle  string
	StrToPreMap map[string]field.PredicateType
}

var plural = pluralize.NewClient()

func NewProto(name string, path string, fields []field.Field) Proto {
	ps := strings.Split(path, "/")
	onlyPath := strings.Join(ps[0:len(ps)-1], "/")
	pkgName := strings.ReplaceAll(onlyPath, "/", ".")

	return Proto{
		Name:        plural.Singular(strings.ToUpper(name[0:1]) + name[1:]),
		OutPath:     path,
		Package:     pkgName,
		GoPackage:   goPackage(onlyPath),
		JavaPackage: javaPackage(pkgName),
		Fields:      fields,
		GenGrpc:     false, // TODO
		GenHttp:     false, // TODO
		primaryKey:  "id",  // TODO
		fieldStyle:  "camel",
		StrToPreMap: field.StrToPreMap,
	}
}

func ModName() string {
	modBytes, err := os.ReadFile("go.mod")
	if err != nil {
		if modBytes, err = os.ReadFile("../go.mod"); err != nil {
			return ""
		}
	}
	return modfile.ModulePath(modBytes)
}

func goPackage(path string) string {
	s := strings.Split(path, "/")
	return ModName() + "/" + path + ";" + s[len(s)-1]
}

func javaPackage(name string) string {
	return name
}

func (p Proto) CreateForm() []field.Field {
	return util.FilterSlice(p.Fields, func(f field.Field) bool {
		return f.Name != p.primaryKey
	})
}

func (p Proto) UpdateForm() []field.Field {
	fs := make([]field.Field, 0, len(p.Fields))

	fs = append(fs, p.PrimaryField())
	fs = append(fs, util.FilterSlice(p.Fields, func(f field.Field) bool {
		return f.Name != p.primaryKey
	})...)
	return fs
}

func (p Proto) ListParams() []field.Predicate {
	fs := make([]field.Predicate, 0, len(p.Fields))
	for _, f := range p.Fields {
		for _, predicate := range f.Predicates {
			fs = append(fs, field.Predicate{
				Name:      f.Name + predicate.Type.String(),
				FieldType: f.FieldType,
				Type:      predicate.Type,
			})
		}
	}
	return fs
}

func (p Proto) PrimaryField() field.Field {
	idField, ok := util.FindSlice(p.Fields, func(f field.Field) bool {
		return f.Name == p.primaryKey
	})
	if !ok {
		idField = field.Field{
			Name:      p.primaryKey,
			FieldType: "int64", // TODO
			Predicates: []field.Predicate{
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

//go:embed tmpl/proto.tmpl
var protoTmpl string

func (p Proto) Generate() error {
	buf := new(bytes.Buffer)

	funcMap := template.FuncMap{
		"ToLower":    strings.ToLower,
		"ToPlural":   plural.Plural,
		"add":        func(a int, b int) int { return a + b },
		"fieldStyle": func(s string) string { return field.StyleFieldMap[p.fieldStyle](s) },
	}
	tmpl, err := template.New("protoTmpl").Funcs(funcMap).Parse(protoTmpl)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, p)
	if err != nil {
		return err
	}

	out := p.OutPath
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
