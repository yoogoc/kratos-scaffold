package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
	"github.com/yoogoc/kratos-scaffold/project_generator"
)

func newGenerateCmd(proto, biz, data, service *cobra.Command) *cobra.Command {
	var generateCmd = &cobra.Command{
		Use:                "generate [NAME]",
		Aliases:            []string{"g"},
		Short:              "gen proto, biz, data, service",
		Long:               `gen proto, biz, data, service`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := proto.RunE(proto, args); err != nil {
				return err
			}
			if err := biz.RunE(biz, args); err != nil {
				return err
			}
			if err := data.RunE(data, args); err != nil {
				return err
			}
			if err := service.RunE(service, args); err != nil {
				return err
			}
			fmt.Println("running wire...")
			if project_generator.IsProjectTypeSingle() {
				if err := util.Exec("wire", "./cmd"); err != nil {
					fmt.Printf("wire failed, please run manually: wire ./cmd\n")
				}
			} else {
				ns := settings.Namespace
				if ns == "" {
					ns = args[0]
				}
				if err := util.Exec("make", "wire-"+ns); err != nil {
					fmt.Printf("wire failed, please run manually: make wire-%s\n", ns)
				}
			}
			return nil
		},
	}

	return generateCmd
}
