package generator

import (
	"os"

	"github.com/yoogoc/kratos-scaffold/pkg/cli"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/scaffoldconfig"
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
	ApiVersion      string
	DI              string
}

func NewBase(setting *cli.EnvSettings, resetNs bool) Base {
	ns := setting.Namespace
	if resetNs {
		if _, err := os.Stat(setting.AppDirName); err != nil && os.IsNotExist(err) {
			ns = ""
		}
	}

	cfg, _, _ := scaffoldconfig.Find(".")
	return Base{
		Namespace:       ns,
		AppDirName:      setting.AppDirName,
		ApiDirName:      setting.ApiDirName,
		PrimaryKey:      setting.PrimaryKey,
		ApiVersion:      setting.GetApiVersion(),
		StrToPreMap:     field.StrToPreMap,
		MaybeGoPackages: field.MaybeGoPackages,
		DI:              cfg.DI,
	}
}
