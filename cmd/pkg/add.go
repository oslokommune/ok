package pkg

import (
	"errors"
	"os"

	"github.com/oslokommune/ok/pkg/pkg/add"
	"github.com/spf13/cobra"
)

var AddCommand = &cobra.Command{
	Use:   "add template [outputFolder]",
	Short: "Add the Boilerplate template to the package manifest with an optional output folder",
	Long: `Add the Boilerplate template to the package manifest with an optional output folder.
The template version is fetched from the latest GitHub release in the template repository.
The output folder is useful when you need multiple instances of the same template with different configurations, for example having multiple instances of the application template.`,
	Example: `ok pkg add databases my-postgres-database
ok pkg add app ecommerce-website
ok pkg add app ecommerce-api
	`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := getArg(args, 0, "")
		outputFolder := getArg(args, 1, templateName)

		result, err := add.Run(PackagesManifestFilename, templateName, outputFolder)
		if err != nil {
			return err
		}

		cmd.Printf("%s (%s) added to %s with output folder name %s\n", result.TemplateName, result.TemplateVersion, PackagesManifestFilename, result.OutputFolder)
		nonExistingConfigFiles := findNonExistingConfigurationFiles(result.VarFiles)
		if len(nonExistingConfigFiles) > 0 {
			cmd.Printf("\nCreate the following configuration files:\n")
			for _, configFile := range nonExistingConfigFiles {
				cmd.Printf("- %s\n", configFile)
			}
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

func findNonExistingConfigurationFiles(varFiles []string) []string {
	var nonExisting []string
	for _, varFile := range varFiles {
		_, err := os.Stat(varFile)
		notExists := errors.Is(err, os.ErrNotExist)
		if notExists {
			nonExisting = append(nonExisting, varFile)
		}
	}
	return nonExisting
}
