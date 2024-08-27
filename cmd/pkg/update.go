package pkg

import (
	"fmt"

	"github.com/oslokommune/ok/pkg/pkg/update"
	"github.com/spf13/cobra"
)

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

		err := update.Run(PackagesManifestFilename, packageName)
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
