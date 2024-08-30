package pkg

import (
	"encoding/json"
	"fmt"

	"github.com/oslokommune/ok/pkg/pkg/config"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
	"github.com/spf13/cobra"
)

var SchemaCommand = &cobra.Command{
	Use: "schema",
}

var SchemaDownloadCommand = &cobra.Command{
	Use: "download",
	RunE: func(cmd *cobra.Command, args []string) error {
		gh, err := githubreleases.GetGitHubClient()
		if err != nil {
			return fmt.Errorf("getting GitHub client: %w", err)
		}
		releases, err := githubreleases.GetLatestReleases()
		if err != nil {
			return fmt.Errorf("getting latest releases: %w", err)
		}

		templateName := args[0]
		templateVersion := releases[templateName]
		githubRef := fmt.Sprintf("%s-%s", templateName, templateVersion)
		templatePath := config.JoinPath("boilerplate/terraform", templateName)
		fileDownloader := githubreleases.NewFileDownloader(gh, boilerplateRepoOwner, boilerplateRepoName, githubRef)
		stacks, err := config.DownloadBoilerplateStacksWithDependencies(cmd.Context(), fileDownloader, templatePath)
		if err != nil {
			return fmt.Errorf("downloading boilerplate stacks: %w", err)
		}
		if len(stacks) == 0 {
			return fmt.Errorf("no stacks found")
		}
		moduleVariables := config.BuildModuleVariables("", stacks[0], stacks, "some/output/folder")
		schema, err := config.TransformModulesToJsonSchema(moduleVariables)
		if err != nil {
			return fmt.Errorf("transforming modules to json schema: %w", err)
		}

		bts, _ := json.MarshalIndent(schema, "", "  ")
		fmt.Println(string(bts))

		return err
	},
}

func init() {
	SchemaCommand.AddCommand(SchemaDownloadCommand)
}
