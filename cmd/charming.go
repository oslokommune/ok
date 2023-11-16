package cmd

import (
	"github.com/oslokommune/ok/charming"
	"github.com/spf13/cobra"
)

func newCharmingCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "charming",
		Short: "Run charming",
		Run: func(cmd *cobra.Command, args []string) {
			charming.Hehe()
		},
	}
}
