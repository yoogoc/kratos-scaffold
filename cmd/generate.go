package cmd

import (
	"github.com/spf13/cobra"
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
			return nil
		},
	}

	return generateCmd
}
