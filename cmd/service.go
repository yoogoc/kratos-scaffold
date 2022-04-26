package cmd

import (
	"github.com/YoogoC/kratos-scaffold/generator"
	"github.com/YoogoC/kratos-scaffold/pkg/field"

	"github.com/spf13/cobra"
)

func newServiceCmd() *cobra.Command {
	var serviceCmd = &cobra.Command{
		Use:                "service [NAME] [FIELD]...",
		Short:              "generate service and service to biz file",
		Long:               `kratos-scaffold service -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Run: func(cmd *cobra.Command, args []string) {
			modelName := args[0]
			service := generator.NewService(modelName, namespace, field.ParseFields(args[1:]))
			err := service.Generate()
			if err != nil {
				panic(err)
			}
		},
	}

	addServiceFlags(serviceCmd)

	return serviceCmd
}

func addServiceFlags(serviceCmd *cobra.Command) {
	serviceCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "app namespace")
}
