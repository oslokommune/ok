package cmd

import (
	"github.com/oslokommune/ok/internal/scriptrunner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envCommand)
}

var envCommand = &cobra.Command{
	Use:   "env",
	Short: "Creates a new `env.yml` file with placeholder values.",
	Run: func(cmd *cobra.Command, args []string) {
		fullArgs := append([]string{"env"}, args...)
		scriptrunner.RunScript("ok.sh", fullArgs)
	},
}
