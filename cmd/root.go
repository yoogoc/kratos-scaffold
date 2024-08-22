package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yoogoc/kratos-scaffold/pkg/cli"
)

var settings = cli.New()

const desc = `kratos-scaffold is a kratos-layout style scaffold.
`

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
		Long:               desc,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	flags := rootCmd.PersistentFlags()

	settings.AddFlags(flags)

	_ = flags.Parse(args)

	serviceCmd := newServiceCmd(args)
	biz := newBizCmd(args)
	protoCmd := newProtoCmd(args)
	dataCmd := newDataCmd(args)

	rootCmd.AddCommand(
		newNewCmd(),
		serviceCmd,
		biz,
		protoCmd,
		dataCmd,
		newGenerateCmd(protoCmd, biz, dataCmd, serviceCmd),
	)

	return rootCmd
}
