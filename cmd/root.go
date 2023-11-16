package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "ok",
		Short: "Command Runner is a simple tool to run a script with subcommands",
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
