package cmd

import (
	"github.com/YoogoC/kratos-scaffold/generator"
	
	"github.com/spf13/cobra"
)

var (
	namespace string
)

func newBizCmd() *cobra.Command {
	var bizCmd = &cobra.Command{
		Use:                "biz [NAME] [FIELD]...",
		Short:              "generate biz file",
		Long:               `kratos-scaffold biz -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Run: func(cmd *cobra.Command, args []string) {
			modelName := args[0]
			biz := generator.NewBiz(modelName, namespace, parseFields(args[1:]))
			err := biz.Generate()
			if err != nil {
				panic(err)
			}
		},
	}
	
	addBizFlags(bizCmd)
	
	return bizCmd
}

func addBizFlags(bizCmd *cobra.Command) {
	bizCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "app namespace")
}
