package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/yoogoc/kratos-scaffold/generator"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"

	"github.com/spf13/cobra"
)

func newBizCmd() *cobra.Command {
	biz := generator.NewBiz(settings)
	var bizCmd = &cobra.Command{
		Use:                "biz [NAME] [FIELD]...",
		Short:              "generate biz file",
		Long:               `kratos-scaffold biz -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBiz(biz, args)
		},
	}

	addBizFlags(bizCmd, biz)

	return bizCmd
}

func addBizFlags(bizCmd *cobra.Command, biz *generator.Biz) {

}

func runBiz(biz *generator.Biz, args []string) error {
	modelName := args[0]

	biz.Name = util.Singular(strcase.ToCamel(modelName))
	if fs, err := field.ParseFields(args[1:]); err != nil {
		return err
	} else {
		biz.Fields = fs
	}

	err := biz.Generate()
	return err
}
