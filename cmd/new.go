package cmd

import (
	pg "github.com/yoogoc/kratos-scaffold/project_generator"

	"github.com/spf13/cobra"
)

var (
	isMono bool
)

func newNewCmd() *cobra.Command {
	project := pg.NewProject()
	var newCmd = &cobra.Command{
		Use:                "new",
		Short:              "generate a new project",
		Long:               `kratos-scaffold new beer-shop`,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			project.Name = args[0]
			project.SetProjectType(isMono)
			return project.Gen()
		},
	}
	addNewFlags(newCmd, project) // inject config struct

	return newCmd
}

func addNewFlags(newCmd *cobra.Command, project *pg.Project) {
	newCmd.PersistentFlags().BoolVarP(&isMono, "mono", "", false, "is mono parent repo?")
	newCmd.PersistentFlags().BoolVarP(&project.IsBff, "bff", "", false, "is bff repo?")
}
