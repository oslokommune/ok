package cmd

import (
	"github.com/oslokommune/ok/charming"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(charmingCommand)
}

var charmingCommand = &cobra.Command{
	Use:   "charming",
	Short: "Run charming",
	Run: func(cmd *cobra.Command, args []string) {
		charming.Hehe()
	},
}
