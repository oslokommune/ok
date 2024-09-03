package pkg

import (
	"fmt"

	"github.com/oslokommune/ok/pkg/pkg/install/interactive"

	"github.com/oslokommune/ok/pkg/pkg/install"
	"github.com/spf13/cobra"
)

var flagInstallInteractive bool
var flagInstallWithLegacyRenderer bool

func init() {
	InstallCommand.Flags().BoolVarP(&flagInstallInteractive,
		"interactive", "i", false, "Select package(s) to install interactively")
	InstallCommand.Flags().BoolVar(&flagInstallWithLegacyRenderer, "use-legacy-boilerplate", false, "Use legacy boilerplate renderer")
}

var InstallCommand = &cobra.Command{
	Use:   "install [outputFolder ...]",
	Short: "Install or update Boilerplate packages.",
	Long: `Install or update Boilerplate packages.

If no arguments are used, the command installs all the packages specified in the package manifest file.

If one or more output folders are specified, the command installs only the packages whose OutputFolder matches the specified folders. (OutputFolder is a field in the package manifest file.)

Set the environment variable BASE_URL to specify where package templates are downloaded from.
`,
	Example: `ok pkg install networking
ok pkg install networking my-app
BASE_URL=../boilerplate/terraform ok pkg install networking my-app
`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, outputFolders []string) error {
		if flagInstallInteractive && len(outputFolders) > 0 {
			return fmt.Errorf("cannot use both --interactive and outputFolder arguments")
		}

		if flagInstallInteractive {
			selectedOutputFolders, err := interactive.SelectPackagesToInstall(PackagesManifestFilename)
			if err != nil {
				return fmt.Errorf("selecting package: %w", err)
			}

			if len(selectedOutputFolders) == 0 {
				fmt.Println("No packages selected. Remember to use space (or x) to select package(s) to install.")
				return nil
			}

			outputFolders = selectedOutputFolders
		}

		err := install.Run(PackagesManifestFilename, outputFolders, flagInstallWithLegacyRenderer)
		if err != nil {
			return fmt.Errorf("installing packages: %w", err)
		}

		return nil
	},
}
