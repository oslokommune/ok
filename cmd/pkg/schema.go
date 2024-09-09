package pkg

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/oslokommune/ok/pkg/pkg/config"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/spf13/cobra"
)

var SchemaCommand = &cobra.Command{
	Use: "schema",
}

var SchemaDownloadCommand = &cobra.Command{
	Use:   "download TEMPLATE [OUTPUT]",
	Short: "Download the JsonSchema for a template",
	RunE: func(cmd *cobra.Command, args []string) error {
		gh, err := githubreleases.GetGitHubClient()
		if err != nil {
			return fmt.Errorf("getting GitHub client: %w", err)
		}

		releases, err := githubreleases.GetLatestReleases()
		if err != nil {
			return fmt.Errorf("getting latest releases: %w", err)
		}

		if len(args) < 1 {
			return fmt.Errorf("missing template name")
		}
		templateName := args[0]
		templateVersion := releases[templateName]
		githubRef := fmt.Sprintf("%s-%s", templateName, templateVersion)

		templatePath := githubreleases.GetTemplatePath(templateName)
		fileDownloader := githubreleases.NewFileDownloader(gh, boilerplateRepoOwner, boilerplateRepoName, githubRef)
		stacks, err := config.DownloadBoilerplateStacksWithDependencies(cmd.Context(), fileDownloader, templatePath)
		if err != nil {
			return fmt.Errorf("downloading boilerplate stacks: %w", err)
		}
		if len(stacks) == 0 {
			return fmt.Errorf("no stacks found")
		}
		moduleVariables := config.BuildModuleVariables(stacks)
		schemaId := fmt.Sprintf("%s-%s", templatePath, templateVersion)
		schema, err := config.TransformModulesToJsonSchema(schemaId, moduleVariables)
		if err != nil {
			return fmt.Errorf("transforming modules to json schema: %w", err)
		}
		bts, _ := json.MarshalIndent(schema, "", "  ")

		outputFile := os.Stdout
		if len(args) > 1 {
			var err error
			outputFile, err = os.Create(args[1])
			if err != nil {
				return fmt.Errorf("creating file: %w", err)
			}
			defer outputFile.Close()
		}

		_, err = fmt.Fprintf(outputFile, "%s\n", bts)
		if err != nil {
			return fmt.Errorf("writing schema to file: %w", err)
		}
		slog.Info(fmt.Sprintf("Schema for %s-%s written to %s\n", templateName, templateVersion, outputFile.Name()))
		return nil
	},
}

func init() {
	SchemaCommand.AddCommand(SchemaDownloadCommand)
}
