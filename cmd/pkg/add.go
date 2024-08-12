package pkg

import (
	"github.com/oslokommune/ok/pkg/pkg/add"
	"github.com/spf13/cobra"
)

var addOutputFolderName string

var AddCommand = &cobra.Command{
	Use:     "add",
	Short:   "Add Boilerplate package to manifest.",
	Example: `ok pkg add <package> [app-name] --output-folder <output-folder>`,
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := args[0]
		appName := ""
		if len(args) >= 2 {
			appName = args[1]
		}

		err := add.Run(PackagesManifestFilename, templateName, addOutputFolderName, appName)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	AddCommand.Flags().StringVar(&addOutputFolderName, "output-folder", "", "Output folder for the new package.")
}
