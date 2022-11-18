package generator

import (
	"os"

	"github.com/yoogoc/kratos-scaffold/pkg/cli"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
)

type Base struct {
	Name            string
	Namespace       string
	AppDirName      string
	ApiDirName      string
	Fields          field.Fields
	StrToPreMap     map[string]field.PredicateType
	MaybeGoPackages []string
	PrimaryKey      string
}

func NewBase(setting *cli.EnvSettings, resetNs bool) Base {
	ns := setting.Namespace
	if resetNs {
		if _, err := os.Stat(setting.AppDirName); err != nil && os.IsNotExist(err) {
			ns = ""
		}
	}

	return Base{
		Namespace:       ns,
		AppDirName:      setting.AppDirName,
		ApiDirName:      setting.ApiDirName,
		PrimaryKey:      setting.PrimaryKey,
		StrToPreMap:     field.StrToPreMap,
		MaybeGoPackages: field.MaybeGoPackages,
	}
}
