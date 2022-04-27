package cmd

import (
	"path"
	"strings"

	"github.com/YoogoC/kratos-scaffold/generator"
	"github.com/YoogoC/kratos-scaffold/pkg/field"
	"github.com/YoogoC/kratos-scaffold/pkg/util"
	"github.com/iancoleman/strcase"

	"github.com/spf13/cobra"
)

func newServiceCmd() *cobra.Command {
	service := generator.NewService(settings)
	var serviceCmd = &cobra.Command{
		Use:                "service [NAME] [FIELD]...",
		Short:              "generate service and service to biz file",
		Long:               `kratos-scaffold service -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runService(service, args)
		},
	}

	addServiceFlags(serviceCmd, service)

	return serviceCmd
}

func addServiceFlags(serviceCmd *cobra.Command, service *generator.Service) {

}

func runService(service *generator.Service, args []string) error {
	modelName := args[0]

	service.Name = util.Singular(strcase.ToCamel(modelName))
	service.Fields = field.ParseFields(args[1:])

	apiModelName := ""
	if service.Namespace != "" {
		apiModelName = path.Join(service.ApiDirName, service.Namespace)
	} else {
		apiModelName = path.Join(service.ApiDirName, strings.ToLower(service.Name))
	}
	// TODO should can set value by flags
	service.ApiPath = path.Join(util.ModName(), apiModelName, "v1")

	err := service.Generate()
	return err
}
