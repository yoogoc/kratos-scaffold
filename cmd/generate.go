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
			if err := proto.RunE(cmd, args); err != nil {
				return err
			}
			if err := biz.RunE(cmd, args); err != nil {
				return err
			}
			if err := data.RunE(cmd, args); err != nil {
				return err
			}
			if err := service.RunE(cmd, args); err != nil {
				return err
			}
			return nil
		},
	}

	return protoCmd
}
