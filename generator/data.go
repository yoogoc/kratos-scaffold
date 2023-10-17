package generator

import (
	"os"
	"path"

	"github.com/yoogoc/kratos-scaffold/pkg/cli"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

type Data struct {
	Base
	NeedAuditField bool
	OrmType        string
	TargetModel    string
}

func NewData(setting *cli.EnvSettings) *Data {
	return &Data{
		Base: NewBase(setting, true),
	}
}

func (d *Data) CreateFields() []*field.Field {
	return d.Fields.CreateFields(d.PrimaryKey)
}

func (d *Data) SoftDelete() bool {
	return d.Fields.HasField("deleted_at")
}

func (d *Data) UpdateFields() []*field.Field {
	return d.Fields.UpdateFields(d.Fields.PrimaryField(d.PrimaryKey))
}

func (d *Data) ParamFields() []*field.Predicate {
	return d.Fields.ParamFields()
}

func (d *Data) InternalPath() string {
	if d.Namespace != "" {
		return path.Join(d.AppDirName, d.Namespace)
	}
	return ""
}

func (d *Data) OutPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(wd, d.InternalPath(), "internal/data")
}

func (d *Data) CurrentPkgPath() string {
	return path.Join(util.ModName(), d.InternalPath(), "internal")
}

func (d *Data) ProtoPkgPath() string {
	return path.Join(util.ModName(), d.ApiDirName, d.TargetModel, "v1")
}
