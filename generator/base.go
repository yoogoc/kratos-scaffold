package generator

import (
	"os"

	"github.com/YoogoC/kratos-scaffold/pkg/cli"
	"github.com/YoogoC/kratos-scaffold/pkg/field"
)

type Base struct {
	Name            string
	Namespace       string
	AppDirName      string
	ApiDirName      string
	Fields          field.Fields
	StrToPreMap     map[string]field.PredicateType
	MaybeGoPackages []string
	primaryKey      string
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
		primaryKey:      setting.PrimaryKey,
		StrToPreMap:     field.StrToPreMap,
		MaybeGoPackages: field.MaybeGoPackages,
	}
}
