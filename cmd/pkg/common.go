package pkg

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

const FlagInteractiveName = "interactive"
const FlagInteractiveShorthand = "i"
const FlagInteractiveUsage = "Select package(s) to install interactively"

const InstallUpdateArgumentDescription = `If no arguments are used, the command installs all the packages specified in the package manifest file.

If one or more output folders are specified, the command installs only the packages whose OutputFolder matches the specified folders. (OutputFolder is a field in the package manifest file.)

Set the environment variable BASE_URL to specify where package templates are downloaded from.
`

var flagInteractive bool
