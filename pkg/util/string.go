package util

import (
	"github.com/gertd/go-pluralize"
)

var plural = pluralize.NewClient()

func Singular(s string) string {
	return plural.Singular(s)
}

func Plural(s string) string {
	return plural.Plural(s)
}
