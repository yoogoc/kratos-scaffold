package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/YoogoC/kratos-scaffold/pkg/util"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"golang.org/x/mod/modfile"
)

type PredicateType int

const (
	PredicateTypeEq = iota
	PredicateTypeContains
	PredicateTypeGt
	PredicateTypeGtE
	PredicateTypeLt
	PredicateTypeLtE
	PredicateTypeIn
)

var (
	presMap = map[PredicateType]string{
		PredicateTypeEq:       "Eq",
		PredicateTypeContains: "Contains",
		PredicateTypeGt:       "Gt",
		PredicateTypeGtE:      "GtEq",
		PredicateTypeLt:       "Lt",
		PredicateTypeLtE:      "LtEq",
		PredicateTypeIn:       "In",
	}
	presEntMap = map[PredicateType]string{
		PredicateTypeEq:       "EQ",
		PredicateTypeContains: "Contains",
		PredicateTypeGt:       "GT",
		PredicateTypeGtE:      "GTE",
		PredicateTypeLt:       "LT",
		PredicateTypeLtE:      "LTE",
		PredicateTypeIn:       "In",
	}
	strToPreMap = map[string]PredicateType{
		"eq":       PredicateTypeEq,
		"contains": PredicateTypeContains,
		"gt":       PredicateTypeGt,
		"gte":      PredicateTypeGtE,
		"lt":       PredicateTypeLt,
		"lte":      PredicateTypeLtE,
		"in":       PredicateTypeIn,
	}

	fieldStyleMap = map[string]func(string) string{
		"snake": strcase.ToSnake,
		"camel": strcase.ToLowerCamel,
	}
)

func (pred PredicateType) String() string {
	return presMap[pred]
}

func (pred PredicateType) EntString() string {
	return presEntMap[pred]
}

func NewPredicateType(s string) PredicateType {
	predicateType, ok := strToPreMap[strings.ToLower(s)]
	if !ok {
		panic(fmt.Sprintf("unknown PredicateType: %s", s))
	}
	return predicateType
}

type Predicate struct {
	Name      string
	Type      PredicateType
	FieldType string
	EntName   string
}

type Field struct {
	Name       string
	FieldType  string
	Predicates []Predicate
}

type Proto struct {
	Name        string
	OutPath     string
	Package     string
	GoPackage   string
	JavaPackage string
	Fields      []Field
	GenGrpc     bool
	GenHttp     bool
	primaryKey  string
	fieldStyle  string
	StrToPreMap map[string]PredicateType
}

var plural = pluralize.NewClient()

func NewProto(name string, path string, fields []Field) Proto {
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
		StrToPreMap: strToPreMap,
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

func (p Proto) CreateForm() []Field {
	return util.FilterSlice(p.Fields, func(f Field) bool {
		return f.Name != p.primaryKey
	})
}

func (p Proto) UpdateForm() []Field {
	fs := make([]Field, 0, len(p.Fields))

	fs = append(fs, p.PrimaryField())
	fs = append(fs, util.FilterSlice(p.Fields, func(f Field) bool {
		return f.Name != p.primaryKey
	})...)
	return fs
}

func (p Proto) ListParams() []Predicate {
	fs := make([]Predicate, 0, len(p.Fields))
	for _, field := range p.Fields {
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

func (p Proto) PrimaryField() Field {
	idField, ok := util.FindSlice(p.Fields, func(f Field) bool {
		return f.Name == p.primaryKey
	})
	if !ok {
		idField = Field{
			Name:      p.primaryKey,
			FieldType: "int64", // TODO
			Predicates: []Predicate{
				{
					Name:      p.primaryKey,
					Type:      PredicateTypeEq,
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
		"fieldStyle": func(s string) string { return fieldStyleMap[p.fieldStyle](s) },
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
