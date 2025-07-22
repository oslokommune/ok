package cmd

import (
	_ "embed"
	"fmt"

	"github.com/oslokommune/ok/pkg/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the `ok` tool and the current latest version available.",
	Run: func(cmd *cobra.Command, args []string) {
		versionInfo, err := version.GetVersionInfo()
		if err != nil {
			fmt.Printf("Error getting version info: %v\n", err)
			return
		}
		fmt.Println(versionInfo)
	},
}
