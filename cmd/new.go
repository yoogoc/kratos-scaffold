package cmd

import (
	"github.com/YoogoC/kratos-scaffold/project_generator"

	"github.com/spf13/cobra"
)

var (
	isMono bool
)

func newNewCmd() *cobra.Command {
	var newCmd = &cobra.Command{
		Use:                "new",
		Short:              "generate a new project",
		Long:               `kratos-scaffold new beer-shop`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Run: func(cmd *cobra.Command, args []string) {
			if isMono {
				err := project_generator.GenMono(args[0])
				if err != nil {
					panic(err)
				}
				return
			}
			err := project_generator.Gen(args[0])
			if err != nil {
				panic(err)
			}
		},
	}
	addNewFlags(newCmd) // inject config struct

	return newCmd
}

func addNewFlags(newCmd *cobra.Command) {
	newCmd.PersistentFlags().BoolVarP(&isMono, "mono", "", false, "is mono parent repo?")
}
