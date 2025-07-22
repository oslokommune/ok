package version

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

// Data holds the version information
var Data Version

// Version represents version information
type Version struct {
	Version string
	Date    string
	Commit  string
}

// GetVersionInfo returns formatted version information similar to the ok version command
func GetVersionInfo() (string, error) {
	result := fmt.Sprintf("Version: %s\nDate:    %s\nCommit:  %s",
		Data.Version,
		Data.Date,
		Data.Commit,
	)

	if Data.Version != "dev" {
		latestVersion, err := githubreleases.GetLatestOkVersion()
		if err != nil {
			result += fmt.Sprintf("\nError getting latest version: %v", err)
		} else {
			currentVersion, err := semver.NewVersion(Data.Version)
			if err != nil {
				result += fmt.Sprintf("\nError parsing version string '%s': %v", Data.Version, err)
			} else {
				if currentVersion.LessThan(latestVersion) {
					result += fmt.Sprintf("\n\nA new release of ok is available: %s → %s", currentVersion, latestVersion)
					result += "\n\nTo update, run the following commands:"
					result += "\n\nbrew update"
					result += "\nbrew upgrade ok"
					result += "\n\nFor other update methods, see https://km.oslo.systems/setup/before-you-start/tools/ok/#updating"
				} else if currentVersion.Equal(latestVersion) {
					result += "\n\nYou are using the latest version."
				}
			}
		}
	}

	return result, nil
}