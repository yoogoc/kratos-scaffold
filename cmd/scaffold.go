package cmd

import (
	"fmt"
	"os"

	"github.com/YoogoC/kratos-scaffold/pkg/cli"
	"github.com/spf13/cobra"
)

var settings = cli.New()

func Execute() {
	rootCmd := newRootCmd(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd(args []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:                "kratos-scaffold",
		Short:              "kratos-scaffold is a kratos-layout style scaffold.",
		Long:               `kratos-scaffold is a kratos-layout style scaffold.`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	flags := rootCmd.PersistentFlags()

	settings.AddFlags(flags)

	_ = flags.Parse(args)

	rootCmd.AddCommand(
		newNewCmd(),
		newServiceCmd(),
		newBizCmd(),
		newProtoCmd(),
		newDataCmd(),
	)

	return rootCmd
}
