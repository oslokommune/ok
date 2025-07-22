package pkg

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"

	"github.com/oslokommune/ok/pkg/pkg/add"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/spf13/cobra"
)

func NewAddCommand() *cobra.Command {
	var flagAddCommandNoSchema bool
	var flagAddCommandNoVarFile bool
	var flagAddCommandVarFile string

	adder := add.NewAdder()

	cmd := &cobra.Command{
		Use:   "add template [outputFolder]",
		Short: "Add the Boilerplate template to the package manifest with an optional output folder",
		Long: `Add the Boilerplate template to the package manifest with an optional output folder.
The template version is fetched from the latest GitHub release in the template repository.
The output folder is useful when you need multiple instances of the same template with different configurations, for example having multiple instances of the application template.`,
		Example: `ok pkg add databases my-postgres-database
ok pkg add app ecommerce-website
ok pkg add app ecommerce-api
	`,
		ValidArgsFunction: addTabCompletion,
		Args:              cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := getArg(args, 0, "")
			outputFolder := getArg(args, 1, templateName)
			currentDir, err := os.Getwd()
			if err != nil {
				cmd.PrintErrf("failed to get current dir: %s\n", err)
				return nil
			}
			consolidatedPackageStructure, err := common.UseConsolidatedPackageStructure(currentDir)
			if err != nil {
				cmd.PrintErrf("failed to check if we should be using consolidated package structure: %s\n", err)
				return nil
			}
			var packagesManifestFilename = common.PackagesManifestFilename
			if !consolidatedPackageStructure {
				packagesManifestFilename = filepath.Join(outputFolder, packagesManifestFilename)
			}

			addSchema := !flagAddCommandNoSchema

			result, err := adder.Run(packagesManifestFilename, templateName, outputFolder, addSchema, consolidatedPackageStructure)
			if err != nil {
				return err
			}

			slog.Info(fmt.Sprintf("%s (%s) added to %s with output folder name %s\n", result.TemplateName, result.TemplateVersion, packagesManifestFilename, result.OutputFolder))
			if consolidatedPackageStructure {
				nonExistingConfigFiles := findNonExistingConfigurationFiles(result.VarFiles)
				if len(nonExistingConfigFiles) > 0 {
					slog.Info("\nCreate the following configuration files:\n")
					for _, configFile := range nonExistingConfigFiles {
						slog.Info(fmt.Sprintf("- %s\n", configFile))
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&flagAddCommandNoSchema, "no-schema", false, "Do not add the JSON schema for the package ")
	cmd.Flags().StringVarP(&flagAddCommandVarFile, "var-file", "v", "default", "Download a var file for the package with the specified name.")
	cmd.Flags().BoolVarP(&flagAddCommandNoVarFile, "no-var-file", "s", false, "Do not download a var file for the package")

	return cmd
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

func addTabCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Only complete the first argument as the second argument is a self-assigned output folder
	if len(args) == 0 {
		return addTabCompletionApp(cmd, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func addTabCompletionApp(cmd *cobra.Command, toComplete string) ([]string, cobra.ShellCompDirective) {
	latest, err := githubreleases.GetLatestReleases()
	if err != nil {
		cmd.PrintErrf("failed to load package manifest: %s\n", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for template := range latest {
		if strings.HasPrefix(template, toComplete) {
			completions = append(completions, template)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
