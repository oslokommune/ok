package cmd

import (
	"github.com/oslokommune/ok/internal/charming"
	"github.com/spf13/cobra"
)

var charmingCommand = &cobra.Command{
	Use:   "charming",
	Short: "Run charming",
	Run: func(cmd *cobra.Command, args []string) {
		charming.Hehe()
	},
}
