package field

import (
	"strings"

	"github.com/YoogoC/kratos-scaffold/pkg/util"
	"github.com/iancoleman/strcase"
)

const DefaultStyleField = "low-camel"

var (
	StyleFields   = []string{"snake", "camel", "low-camel"}
	StyleFieldMap = map[string]func(string) string{
		"snake":     strcase.ToSnake,
		"low-camel": strcase.ToLowerCamel,
		"camel":     strcase.ToCamel,
	}
)

type Fields []*Field

type Field struct {
	Name       string
	FieldType  string
	Predicates []*Predicate
}

func (fs Fields) CreateFields(primaryKey string) []*Field {
	return util.FilterSlice(fs, func(f *Field) bool {
		return f.Name != primaryKey
	})
}

func (fs Fields) UpdateFields(primaryField *Field) []*Field {
	result := make([]*Field, 0, len(fs))
	result = append(result, primaryField)

	return append(result, util.FilterSlice(fs, func(f *Field) bool {
		return f.Name != primaryField.Name
	})...)
}

func (fs Fields) ParamFields() []*Predicate {
	result := make([]*Predicate, 0, len(fs))
	for _, f := range fs {
		for _, predicate := range f.Predicates {
			result = append(result, &Predicate{
				Name:      f.Name + predicate.Type.String(),
				FieldType: f.FieldType,
				Type:      predicate.Type,
				EntName:   EntName(f.Name) + predicate.Type.EntString(),
			})
		}
	}
	return result
}

func EntName(s string) string {
	s = strcase.ToCamel(s)
	if len(s) < 2 || strings.ToLower(s[len(s)-2:]) != "id" {
		return s
	}
	return s[:len(s)-2] + "ID"
}

func ParseFields(strs []string) []*Field {
	var fs []*Field
	for _, str := range strs {
		ss := strings.Split(str, ":")
		var pres []*Predicate
		if len(ss) > 2 {
			for _, p := range strings.Split(ss[2], ",") {
				pres = append(pres, &Predicate{
					Name:      ss[0],
					Type:      NewPredicateType(p),
					FieldType: ss[1],
				})
			}
		}
		fs = append(fs, &Field{
			Name:       ss[0],
			FieldType:  ss[1],
			Predicates: pres,
		})
	}
	return fs
}
