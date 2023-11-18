package cmd

import (
	"os"

	"github.com/oslokommune/ok/internal/charming"
	"github.com/spf13/cobra"
)

func init() {
	if os.Getenv("OK_ENABLE_EXPERIMENTAL") == "true" {
		rootCmd.AddCommand(charmingCommand)
	}
}

var charmingCommand = &cobra.Command{
	Use:   "charming",
	Short: "Run charming",
	Run: func(cmd *cobra.Command, args []string) {
		charming.Hehe()
	},
}
