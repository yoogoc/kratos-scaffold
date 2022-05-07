package cmd

import (
	"fmt"
	"os"

	"github.com/YoogoC/kratos-scaffold/pkg/cli"
	"github.com/spf13/cobra"
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

	service := newServiceCmd()
	biz := newBizCmd()
	proto := newProtoCmd()
	data := newDataCmd()

	rootCmd.AddCommand(
		newNewCmd(),
		service,
		biz,
		proto,
		data,
		newGenerateCmd(proto, biz, data, service),
	)

	return rootCmd
}
