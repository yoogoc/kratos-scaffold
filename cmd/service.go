package cmd

import (
	"path"

	"github.com/iancoleman/strcase"
	"github.com/yoogoc/kratos-scaffold/generator"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
	"github.com/yoogoc/kratos-scaffold/project_generator"

	"github.com/spf13/cobra"
)

func newServiceCmd(args []string) *cobra.Command {
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

	addServiceFlags(serviceCmd, service, args)

	return serviceCmd
}

func addServiceFlags(serviceCmd *cobra.Command, service *generator.Service, args []string) {
	flags := serviceCmd.PersistentFlags()
	flags.ParseErrorsWhitelist.UnknownFlags = true
	flags.StringVarP(&service.ApiPath, "api-path", "", "", "proto path, default is <current-mod-name>/<api-dir-name>/<namespace | name>/v1")
	flags.BoolVarP(&service.GenHttp, "http", "", false, "generate xx.http.pb.go")
	_ = flags.Parse(args)
}

func runService(service *generator.Service, args []string) error {
	modelName := args[0]

	service.Name = util.Singular(strcase.ToCamel(modelName))
	if fs, err := field.ParseFields(args[1:]); err != nil {
		return err
	} else {
		service.Fields = fs
	}

	if service.ApiPath == "" {
		apiModelName := ""
		if project_generator.IsProjectTypeSingle() {
			apiModelName = service.ApiDirName
		} else {
			if service.Namespace != "" {
				apiModelName = path.Join(service.ApiDirName, service.Namespace)
			} else {
				apiModelName = path.Join(service.ApiDirName, modelName)
			}
		}

		service.ApiPath = path.Join(util.ModName(), apiModelName, "v1")
	}

	err := service.Generate()
	return err
}
