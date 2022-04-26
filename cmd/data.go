package cmd

import (
	"github.com/YoogoC/kratos-scaffold/generator"
	"github.com/YoogoC/kratos-scaffold/pkg/field"

	"github.com/spf13/cobra"
)

var (
	ormType string
)

func newDataCmd() *cobra.Command {
	var dataCmd = &cobra.Command{
		Use:                "data [NAME] [FIELD]...",
		Short:              "generate data, data and data to biz file",
		Long:               `kratos-scaffold data -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Run: func(cmd *cobra.Command, args []string) {
			modelName := args[0]
			data := generator.NewDataEnt(modelName, namespace, field.ParseFields(args[1:]))
			err := data.Generate()
			if err != nil {
				panic(err)
			}
		},
	}

	addDataFlags(dataCmd)

	return dataCmd
}

func addDataFlags(dataCmd *cobra.Command) {
	dataCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "app namespace")
	dataCmd.PersistentFlags().StringVar(&ormType, "orm", "ent", "app namespace")
}
