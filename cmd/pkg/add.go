package pkg

import (
	"github.com/oslokommune/ok/pkg/pkg/add"
	"github.com/spf13/cobra"
)

var flagAddOutputFolder string

var AddCommand = &cobra.Command{
	Use:   "add template [stack-name]",
	Short: "Add Boilerplate template to packages manifest with an optional stack name",
	Long: `Add Boilerplate template to packages manifest with an optional stack name.
The template version is fetched from the latest release on GitHub and added to the packages manifest without applying the template.
The stack name is useful to define if you need multiple instances of the same template with different configurations. For example having multiple apps in the same project.`,
	Example: `ok pkg add databases my-postgres-database
ok pkg add app website --output-folder ecommerce
ok pkg add app api --output-folder ecommerce
	`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := getArg(args, 0, "")
		stackName := getArg(args, 1, templateName)

		err := add.Run(PackagesManifestFilename, templateName, flagAddOutputFolder, stackName)
		if err != nil {
			return err
		}

		return nil
	},
}

func getArg(args []string, index int, fallback string) string {
	if len(args) > index {
		return args[index]
	}
	return fallback
}

func init() {
	AddCommand.Flags().StringVarP(&flagAddOutputFolder, "output-folder", "o", "", "Output folder for the new package. Defaults to the current working directory.")
}
