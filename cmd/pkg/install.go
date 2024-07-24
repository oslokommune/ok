package pkg

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/install"
	"github.com/oslokommune/ok/pkg/pkg/install/interactive"
	"github.com/spf13/cobra"
)

var installSelectFlag bool

func init() {
	InstallCommand.Flags().BoolVarP(&installSelectFlag, "interactive", "i", false, "Select package interactively")
}

var InstallCommand = &cobra.Command{
	Use:   "install [outputFolder ...]",
	Short: "Install or update Boilerplate packages.",
	Long: `Install or update Boilerplate packages.

If no arguments are used, the command installs all the packages specified in the package manifest file.

If one or more output folders are specified, the command installs only the packages whose OutputFolder matches the specified folders. (OutputFolder is a field in the package manifest file.)

Set the environment variable BASE_URL to specify where package templates are downloaded from. 
`,
	Example: `ok install networking
ok install networking my-app
BASE_URL=../boilerplate/terraform ok install networking my-app
`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, stacks []string) error {
		if installSelectFlag {
			answer, err := interactive.Run(PackagesManifestFilename)
			if err != nil {
				return fmt.Errorf("selecting package: %w", err)
			}

			if answer.Aborted {
				fmt.Println("Aborted")
				return nil
			}

			stacks = []string{answer.Choice}
		}

		err := install.Run(PackagesManifestFilename, stacks)
		if err != nil {
			return err
		}

		return nil
	},
}
