package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/yoogoc/kratos-scaffold/generator"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"

	"github.com/spf13/cobra"
)

func newDataCmd() *cobra.Command {
	dataEnt := generator.NewData(settings)
	var dataCmd = &cobra.Command{
		Use:                "data [NAME] [FIELD]...",
		Short:              "generate data, data and data to biz file",
		Long:               `kratos-scaffold data -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			if dataEnt.OrmType == "ent" {
				return runDataEnt(dataEnt, args)
			}
			if dataEnt.OrmType == "proto" {
				return runDataProto(dataEnt, args)
			}
			return nil
		},
	}

	addDataFlags(dataCmd, dataEnt)

	return dataCmd
}

func addDataFlags(dataCmd *cobra.Command, dataEnt *generator.Data) {
	dataCmd.PersistentFlags().BoolVarP(&dataEnt.NeedAuditField, "audit-field", "", true, "auto generate created_at and update_at fields, default is true")
	dataCmd.PersistentFlags().StringVarP(&dataEnt.OrmType, "orm-type", "t", "ent", "orm type, value in (ent, proto), default is ent")
	dataCmd.PersistentFlags().StringVarP(&dataEnt.TargetModel, "target-model", "m", "", "proto-orm required")
}

func runDataEnt(dataEnt *generator.Data, args []string) error {
	modelName := args[0]

	dataEnt.Name = util.Singular(strcase.ToCamel(modelName))
	if fs, err := field.ParseFields(args[1:]); err != nil {
		return err
	} else {
		dataEnt.Fields = fs
	}

	if err := dataEnt.GenerateMigration(); err != nil {
		return err
	}

	if err := dataEnt.GenerateEnt(); err != nil {
		return err
	}
	return nil
}

func runDataProto(dataEnt *generator.Data, args []string) error {
	modelName := args[0]

	dataEnt.Name = util.Singular(strcase.ToCamel(modelName))
	if fs, err := field.ParseFields(args[1:]); err != nil {
		return err
	} else {
		dataEnt.Fields = fs
	}

	err := dataEnt.GenerateProto()
	return err
}
