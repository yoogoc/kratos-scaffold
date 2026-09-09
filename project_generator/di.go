package project_generator

import (
	"bytes"
	"fmt"
	"regexp"
)

var activeDI = "wire"

const kitPaginationAliases = `package biz

import "github.com/yoogoc/kratos-kit/pagination"

type OrderBy = pagination.OrderBy
type ListOption = pagination.ListOption
type ListOptions = pagination.ListOptions

var DefaultListOptions = pagination.DefaultListOptions
var ApplyListOptions = pagination.ApplyListOptions
var WithPage = pagination.WithPage
var WithPageSize = pagination.WithPageSize
var WithOrderBy = pagination.WithOrderBy
var FormatOrderBy = pagination.FormatOrderBy
var ParseOrderBy = pagination.ParseOrderBy
`

func isDo() bool { return activeDI == "do" }

func providerDeclaration(packageName string, constructors ...string) string {
	if !isDo() {
		return fmt.Sprintf("package %s\n\nimport \"github.com/google/wire\"\n\nvar ProviderSet = wire.NewSet(%s)\n", packageName, joinConstructors(constructors))
	}
	return fmt.Sprintf("package %s\n", packageName)
}

func joinConstructors(values []string) string {
	out := ""
	for i, value := range values {
		if i > 0 {
			out += ", "
		}
		out += value
	}
	return out
}

var wireSetPattern = regexp.MustCompile(`(?m)^var ProviderSet = wire\.NewSet\([^\n]*\)\n?`)

func stripWire(content []byte) []byte {
	content = bytes.ReplaceAll(content, []byte("\t\"github.com/google/wire\"\n"), nil)
	content = bytes.ReplaceAll(content, []byte("\"github.com/google/wire\"\n"), nil)
	return wireSetPattern.ReplaceAll(content, nil)
}
