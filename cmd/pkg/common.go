package pkg

import (
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/spf13/cobra"
	"path"
)

const FlagInteractiveName = "interactive"
const FlagInteractiveShorthand = "i"
const FlagInteractiveUsage = "Select package(s) to install interactively"

const InstallUpdateArgumentDescription = `If no arguments are used, the command installs all the packages specified in the package manifest file.

If one or more output folders are specified, the command installs only the packages whose OutputFolder matches the specified folders. (OutputFolder is a field in the package manifest file.)

Set the environment variable BASE_URL to specify where package templates are downloaded from.
`

var flagCwd string
var flagInteractive bool

// argsContainsElement checks if an element is in an array which
// is used in the tab completion functions to filter out already selected elements
func argsContainsElement[T comparable](arr []T, element T) bool {
	for _, e := range arr {
		if e == element {
			return true
		}
	}
	return false
}

func AddCwdFlag(cmd *cobra.Command, variable *string) {
	cmd.Flags().StringVarP(variable,
		common.FlagNameCwd,
		"",
		".",
		"Set the current working directory",
	)
}

func GetPackageManifestPath(workingDirectory string) string {
	return path.Join(workingDirectory, "packages.yml")
}
