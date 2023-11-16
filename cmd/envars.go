package cmd

import (
	"github.com/oslokommune/ok/scriptrunner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envarsCommand)
}

var envarsCommand = &cobra.Command{
	Use:   "envars",
	Short: "Exports the values in env.yml as environment variables.",
	Run: func(cmd *cobra.Command, args []string) {
		fullArgs := append([]string{"envars"}, args...)
		scriptrunner.RunScript("ok.sh", fullArgs)
	},
}
