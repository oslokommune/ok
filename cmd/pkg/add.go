package pkg

import (
	"errors"
	"fmt"
	cmdPkgCommon "github.com/oslokommune/ok/cmd/pkg/common"
	"log/slog"
	"os"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/add"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/spf13/cobra"
)

var flagAddCommandUpdateSchema bool

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
	ValidArgsFunction: addTabCompletion,
	Args:              cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := getArg(args, 0, "")
		outputFolder := getArg(args, 1, templateName)

		result, err := add.Run(flagPackageFile, templateName, outputFolder, flagAddCommandUpdateSchema)
		if err != nil {
			return err
		}

		slog.Info(fmt.Sprintf("%s (%s) added to %s with output folder name %s\n", result.TemplateName, result.TemplateVersion, flagPackageFile, result.OutputFolder))
		nonExistingConfigFiles := findNonExistingConfigurationFiles(result.VarFiles)
		if len(nonExistingConfigFiles) > 0 {
			slog.Info("\nCreate the following configuration files:\n")
			for _, configFile := range nonExistingConfigFiles {
				slog.Info(fmt.Sprintf("- %s\n", configFile))
			}
		}
		return nil
	},
}

func init() {
	cmdPkgCommon.AddPackageFileFlag(AddCommand, &flagPackageFile)
	AddCommand.Flags().BoolVar(&flagAddCommandUpdateSchema, "update-schema", true, "Update the JSON schema for affected packages")
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
