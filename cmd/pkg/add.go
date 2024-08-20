package pkg

import (
	"errors"
	"os"

	"github.com/oslokommune/ok/pkg/pkg/add"
	"github.com/spf13/cobra"
)

var AddCommand = &cobra.Command{
	Use:   "add template [outputFolder]",
	Short: "Add Boilerplate template to packages manifest with an optional output folder",
	Long: `Add Boilerplate template to packages manifest with an optional output folder.
The template version is fetched from the latest release on GitHub and added to the packages manifest without applying the template.
The output folder is useful to define if you need multiple instances of the same template with different configurations. For example having multiple apps in the same project.`,
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
			cmd.Printf("\nCreate these following configuration files:\n")
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
