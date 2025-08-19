package pkg

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"os"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/add"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/spf13/cobra"
)

func NewAddCommand(ghReleases add.GitHubReleases) *cobra.Command {
	var flagAddCommandNoSchema bool
	var flagAddCommandNoVarFile bool
	var flagAddCommandVarFile string

	cmd := &cobra.Command{
		Use:   "add <template> [outputFolder]",
		Short: "Add the Boilerplate template to the package manifest with an optional output folder",
		Long: `Add the Boilerplate template to the package manifest with an optional output folder.
The template version is fetched from the latest GitHub release in the template repository.
The output folder is useful when you need multiple instances of the same template with different configurations, for example having multiple instances of the application template.`,
		Example: `ok pkg add databases my-postgres-database
ok pkg add app ecommerce-website
ok pkg add app ecommerce-api
BASE_URL=../boilerplate/terraform ok pkg add networking
	`,
		ValidArgsFunction: addTabCompletion,
		Args:              cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := getArg(args, 0, "")
			outputFolder := getArg(args, 1, templateName)

			if flagAddCommandVarFile != "default" && flagAddCommandNoVarFile {
				return fmt.Errorf("cannot use both --var-file and --%s flags", add.FlagNoVar)
			}

			currentDir, err := os.Getwd()
			if err != nil {
				cmd.PrintErrf("failed to get current dir: %s\n", err)
				return nil
			}

			adder := add.NewAdder(ghReleases)

			err = adder.Run(add.Options{
				BaseUrl:         os.Getenv(common.BaseUrlEnvName),
				CurrentDir:      currentDir,
				TemplateName:    templateName,
				OutputFolder:    outputFolder,
				AddSchema:       !flagAddCommandNoSchema,
				DownloadVarFile: !flagAddCommandNoVarFile,
				VarFile:         flagAddCommandVarFile,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&flagAddCommandNoSchema, "no-schema", false, "Do not add the JSON schema for the package ")
	cmd.Flags().StringVarP(&flagAddCommandVarFile, "var-file", "v", "default", "Download a var file for the package with the specified name")
	cmd.Flags().BoolVarP(&flagAddCommandNoVarFile, add.FlagNoVar, "s", false, "Do not download a var file for the package")

	return cmd
}

func getArg(args []string, index int, fallback string) string {
	if len(args) > index {
		return args[index]
	}
	return fallback
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
