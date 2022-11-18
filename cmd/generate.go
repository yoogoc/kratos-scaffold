package cmd

import (
	"github.com/spf13/cobra"
)

func newGenerateCmd(proto, biz, data, service *cobra.Command) *cobra.Command {
	var protoCmd = &cobra.Command{
		Use:                "generate [NAME]",
		Aliases:            []string{"g"},
		Short:              "gen proto, biz, data, service",
		Long:               `gen proto, biz, data, service`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := proto.Execute(); err != nil {
				return err
			}
			if err := biz.Execute(); err != nil {
				return err
			}
			if err := data.Execute(); err != nil {
				return err
			}
			if err := service.Execute(); err != nil {
				return err
			}
			return nil
		},
	}

	return protoCmd
}
