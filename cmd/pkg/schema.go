package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/schema"
	"log/slog"
	"os"

	"github.com/oslokommune/ok/pkg/pkg/config"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/spf13/cobra"
)

func init() {
	SchemaCommand.AddCommand(SchemaDownloadCommand)
	AddCwdFlag(SchemaDownloadCommand, &flagCwd)
}

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
		manifest, err := common.LoadPackageManifest(flagCwd)
		if err != nil {
			return fmt.Errorf("could not load package manifest: %w", err)
		}
		templateName := args[0]
		templateVersion := releases[templateName]
		githubRef := fmt.Sprintf("%s-%s", templateName, templateVersion)

		templatePath := githubreleases.GetTemplatePath(manifest.PackagePrefix(), templateName)
		fileDownloader := githubreleases.NewFileDownloader(gh, common.BoilerplateRepoOwner, common.BoilerplateRepoName, githubRef)
		stacks, err := config.DownloadBoilerplateStacksWithDependencies(cmd.Context(), fileDownloader, templatePath)
		if err != nil {
			return fmt.Errorf("downloading boilerplate stacks: %w", err)
		}
		if len(stacks) == 0 {
			return fmt.Errorf("no stacks found")
		}
		moduleVariables := schema.BuildModuleVariables(stacks)
		schemaId := fmt.Sprintf("%s-%s", templatePath, templateVersion)
		schema, err := schema.TransformModulesToJsonSchema(schemaId, moduleVariables)
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
