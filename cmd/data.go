package cmd

import (
	"github.com/YoogoC/kratos-scaffold/generator"
	"github.com/YoogoC/kratos-scaffold/pkg/field"
	"github.com/YoogoC/kratos-scaffold/pkg/util"
	"github.com/iancoleman/strcase"

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

	addDataFlags(dataCmd)

	return dataCmd
}

func addDataFlags(dataCmd *cobra.Command) {
	// dataCmd.PersistentFlags().StringVar(&ormType, "orm", "ent", "app namespace")
}

func runDataEnt(dataEnt *generator.DataEnt, args []string) error {
	modelName := args[0]

	dataEnt.Name = util.Singular(strcase.ToCamel(modelName))
	dataEnt.Fields = field.ParseFields(args[1:])

	err := dataEnt.Generate()
	return err
}
