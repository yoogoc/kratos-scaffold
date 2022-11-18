package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/yoogoc/kratos-scaffold/generator"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"

	"github.com/spf13/cobra"
)

func newDataCmd() *cobra.Command {
	dataEnt := generator.NewDataEnt(settings)
	var dataCmd = &cobra.Command{
		Use:                "data [NAME] [FIELD]...",
		Short:              "generate data, data and data to biz file",
		Long:               `kratos-scaffold data -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDataEnt(dataEnt, args)
		},
	}

	addDataFlags(dataCmd, dataEnt)

	return dataCmd
}

func addDataFlags(dataCmd *cobra.Command, dataEnt *generator.DataEnt) {
	dataCmd.PersistentFlags().BoolVarP(&dataEnt.NeedAuditField, "audit-field", "", true, "auto generate created_at and update_at fields, default is true")
}

func runDataEnt(dataEnt *generator.DataEnt, args []string) error {
	modelName := args[0]

	dataEnt.Name = util.Singular(strcase.ToCamel(modelName))
	if fs, err := field.ParseFields(args[1:]); err != nil {
		return err
	} else {
		dataEnt.Fields = fs
	}

	err := dataEnt.Generate()
	return err
}
