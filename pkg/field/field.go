package field

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
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
	FieldType  TypeField
	Predicates []*Predicate
}

func (fs Fields) CreateFields(primaryKey string) []*Field {
	return util.FilterSlice(fs, func(f *Field) bool {
		fn := strings.ToLower(strcase.ToSnake(f.Name))
		return fn != strings.ToLower(strcase.ToSnake(primaryKey)) &&
			fn != "created_at" &&
			fn != "updated_at" &&
			fn != "created_by" &&
			fn != "updated_by" &&
			fn != "deleted_at"
	})
}

func (fs Fields) HasField(field string) bool {
	_, ok := util.FindSlice(fs, func(f *Field) bool {
		return strcase.ToSnake(f.Name) == strcase.ToSnake(field)
	})
	return ok
}

func (fs Fields) UpdateFields(primaryField *Field) []*Field {
	return util.FilterSlice(fs, func(f *Field) bool {
		fn := strings.ToLower(strcase.ToSnake(f.Name))
		return fn != strings.ToLower(strcase.ToSnake(primaryField.Name)) &&
			fn != "created_at" &&
			fn != "updated_at" &&
			fn != "created_by" &&
			fn != "updated_by" &&
			fn != "deleted_at"
	})
}

func (fs Fields) ParamFields() []*Predicate {
	result := make([]*Predicate, 0, len(fs))
	for _, f := range fs {
		for _, predicate := range f.Predicates {
			result = append(result, &Predicate{
				Name:       f.Name + predicate.Type.StringProto(),
				SourceName: f.Name,
				FieldType:  f.FieldType,
				Type:       predicate.Type,
				EntName:    EntName(f.Name) + predicate.Type.EntString(),
			})
		}
	}
	return result
}

func (fs Fields) PrimaryField(primaryKey string) *Field {
	idField, ok := util.FindSlice(fs, func(f *Field) bool {
		return f.Name == primaryKey
	})
	if !ok {
		idField = &Field{
			Name:      primaryKey,
			FieldType: DefaultPrimaryFieldType,
			Predicates: []*Predicate{
				{
					Name:      primaryKey,
					Type:      PredicateTypeEq,
					FieldType: DefaultPrimaryFieldType,
				},
			},
		}
	}
	return idField
}

func EntName(s string) string {
	s = strcase.ToCamel(s)
	if len(s) < 2 || strings.ToLower(s[len(s)-2:]) != "id" {
		return s
	}
	return s[:len(s)-2] + "ID"
}

func ParseFields(strs []string) ([]*Field, error) {
	var fs []*Field
	for _, str := range strs {
		ss := strings.Split(str, ":")
		var pres []*Predicate
		t, ok := types[ss[1]]
		if !ok {
			return nil, errors.New(fmt.Sprintf("unknown type: %s\n", ss[1]))
		}
		if len(ss) > 2 {
			for _, p := range strings.Split(ss[2], ",") {
				pres = append(pres, &Predicate{
					Name:      ss[0],
					Type:      NewPredicateType(p),
					FieldType: t,
				})
			}
		}
		fs = append(fs, &Field{
			Name:       ss[0],
			FieldType:  t,
			Predicates: pres,
		})
	}
	return fs, nil
}
