package cmd

import (
	"github.com/oslokommune/ok/scriptrunner"
	"github.com/spf13/cobra"
)

func newEnvCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "Creates a new env.yml file with placeholder values.",
		Run: func(cmd *cobra.Command, args []string) {
			fullArgs := append([]string{"env"}, args...)
			scriptrunner.RunScript("ok.sh", fullArgs)
		},
	}
}
