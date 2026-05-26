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
	IsArray    bool
	IsNillable bool
	Predicates []*Predicate
}

func (fs Fields) ProtoImports() []string {
	var imports []string
	for _, f := range fs {
		if f.IsArray || f.FieldType.IsJson() {
			imports = append(imports, "google/protobuf/struct.proto")
			break
		}
	}
	return imports
}

func (fs Fields) HasTimeField() bool {
	for _, f := range fs {
		if f.FieldType == TypeTime || f.FieldType == TypeDate {
			return true
		}
	}
	return false
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

func (fs Fields) UpdateProtoFields() []*Field {
	return util.FilterSlice(fs, func(f *Field) bool {
		fn := strings.ToLower(strcase.ToSnake(f.Name))
		return fn != "created_at" &&
			fn != "updated_at" &&
			fn != "created_by" &&
			fn != "updated_by" &&
			fn != "deleted_at"
	})
}

func (fs Fields) TransferFields() []*Field {
	return util.FilterSlice(fs, func(f *Field) bool {
		fn := strings.ToLower(strcase.ToSnake(f.Name))
		return fn != "deleted_at"
	})
}

// EntityFields returns fields for entity struct, excluding deleted_at
func (fs Fields) EntityFields() []*Field {
	return util.FilterSlice(fs, func(f *Field) bool {
		fn := strings.ToLower(strcase.ToSnake(f.Name))
		return fn != "deleted_at"
	})
}

// ParamFields returns predicates for query params, excluding deleted_at
func (fs Fields) ParamFields() []*Predicate {
	result := make([]*Predicate, 0, len(fs))
	for _, f := range fs {
		fn := strings.ToLower(strcase.ToSnake(f.Name))
		if fn == "deleted_at" {
			continue
		}
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
		ts := strings.Split(ss[1], ",")
		t, ok := types[ts[0]]
		if !ok {
			return nil, errors.New(fmt.Sprintf("unknown type: %s\n", ss[1]))
		}
		isArray := contains(ts[1:], "a", "array", "slice")
		isNillable := contains(ts[1:], "n", "nil", "null")

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
			IsArray:    isArray,
			IsNillable: isNillable,
			Predicates: pres,
		})
	}
	return fs, nil
}
