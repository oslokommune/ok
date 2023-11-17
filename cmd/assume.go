package cmd

import (
	"github.com/oslokommune/ok/pkg/toggle"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(assumeCommand)
}

var assumeCommand = &cobra.Command{
	Use:   "assume",
	Short: "Toggle assume_cd_role in app stack",
	Run: func(cmd *cobra.Command, args []string) {
		toggle.Assume()
	},
}
