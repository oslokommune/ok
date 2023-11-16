package cmd

import (
	"github.com/oslokommune/ok/scriptrunner"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of the ok tool and the current latest version available.",
		Run: func(cmd *cobra.Command, args []string) {
			fullArgs := append([]string{"version"}, args...)
			scriptrunner.RunScript("ok.sh", fullArgs)
		},
	}
}
