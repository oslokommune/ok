package cmd

import (
	_ "embed"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
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

		if VersionData.Version != "dev" {
			latestVersion, err := githubreleases.GetLatestOkVersion()
			if err != nil {
				fmt.Printf("Error getting latest version: %v\n", err)
				return
			}

			currentVersion, err := semver.NewVersion(VersionData.Version)
			if err != nil {
				fmt.Printf("Error parsing version string '%s': %v\n", VersionData.Version, err)
				return
			}

			if currentVersion.LessThan(latestVersion) {
				fmt.Printf("\nA new release of ok is available: %s → %s\n", currentVersion, latestVersion)
				fmt.Println("\nTo update, run the following commands:")
				fmt.Println("\nbrew update")
				fmt.Println("brew upgrade ok")
				fmt.Println("\nFor other update methods, see https://km.oslo.systems/setup/before-you-start/tools/ok/#updating")
			} else if currentVersion.Equal(latestVersion) {
				fmt.Println("\nYou are using the latest version.")
			}
		}

	},
}

type Version struct {
	Version string
	Date    string
	Commit  string
}
