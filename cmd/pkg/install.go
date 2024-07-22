package pkg

import (
	"github.com/oslokommune/ok/pkg/pkg/install"
	"github.com/spf13/cobra"
)

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
		err := install.Run(PackagesManifestFilename, stacks)
		if err != nil {
			return err
		}

		return nil
	},
}
