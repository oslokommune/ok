package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "ok",
		Short: "The ok tool.",
		Long:  "The ok tool.",
	}

	rootCmd.AddCommand(newBootstrapCommand(),
		newScaffoldCommand(),
		newEnvCommand(),
		newEnvarsCommand(),
		newGetTemplateCommand(),
		newForwardCommand(),
		newVersionCommand(),
		newAssumeCommand(),
		newCharmingCommand())

	return rootCmd
}
