package field

import "github.com/iancoleman/strcase"

var (
	StyleFieldMap = map[string]func(string) string{
		"snake": strcase.ToSnake,
		"camel": strcase.ToLowerCamel,
	}
)

type Field struct {
	Name       string
	FieldType  string
	Predicates []Predicate
}
