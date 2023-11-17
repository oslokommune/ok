package cmd

import (
	"github.com/oslokommune/ok/internal/scriptrunner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the ok tool and the current latest version available.",
	Run: func(cmd *cobra.Command, args []string) {
		fullArgs := append([]string{"version"}, args...)
		scriptrunner.RunScript("ok.sh", fullArgs)
	},
}
