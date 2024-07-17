package cmd

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
)

var VersionData Version

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the ok tool and the current latest version available.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", VersionData.Version)
		fmt.Printf("Date:    %s\n", VersionData.Date)
		fmt.Printf("Commit:  %s\n", VersionData.Commit)
	},
}

type Version struct {
	Version string
	Date    string
	Commit  string
}
