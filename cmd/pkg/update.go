package pkg

import (
	"fmt"

	"github.com/oslokommune/ok/pkg/pkg/update"
	"github.com/spf13/cobra"
)

var flagUpdateCommandUpdateSchema bool

var UpdateCommand = &cobra.Command{
	Use:   "update [package-name]",
	Short: "Update Boilerplate package manifest",
	Long: `Update Boilerplate package manifest.
If a package name is provided, only that package will be updated.
If no package name is provided, all packages will be updated.`,
	Example: `  ok pkg update
  ok pkg update my-package`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var packageName string
		if len(args) == 1 {
			packageName = args[0]
		}

		err := update.Run(PackagesManifestFilename, packageName, flagUpdateCommandUpdateSchema)
		if err != nil {
			return err
		}

		if packageName != "" {
			fmt.Printf("Updated package: %s\n", packageName)
		} else {
			fmt.Println("Updated all packages")
		}

		return nil
	},
}

func init() {
	UpdateCommand.Flags().BoolVar(&flagUpdateCommandUpdateSchema, "update-schema", true, "Update the JSON schema for affected packages")
}
