package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/yoogoc/kratos-scaffold/generator"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"

	"github.com/spf13/cobra"
)

func newDataCmd(args []string) *cobra.Command {
	data := generator.NewData(settings)
	var dataCmd = &cobra.Command{
		Use:                "data [NAME] [FIELD]...",
		Short:              "generate data, data and data to biz file",
		Long:               `kratos-scaffold data -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			if data.OrmType == "ent" {
				return runDataEnt(data, args)
			}
			if data.OrmType == "proto" {
				return runDataProto(data, args)
			}
			return nil
		},
	}

	addDataFlags(dataCmd, data, args)

	return dataCmd
}

func addDataFlags(dataCmd *cobra.Command, dataEnt *generator.Data, args []string) {
	flags := dataCmd.PersistentFlags()
	flags.ParseErrorsWhitelist.UnknownFlags = true
	flags.BoolVarP(&dataEnt.NeedAuditField, "audit-field", "", true, "auto generate created_at and update_at fields, default is true")
	flags.StringVarP(&dataEnt.OrmType, "orm-type", "t", "ent", "orm type, value in (ent, proto), default is ent")
	flags.StringVarP(&dataEnt.TargetModel, "target-model", "m", "", "proto-orm required")
	_ = flags.Parse(args)
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
