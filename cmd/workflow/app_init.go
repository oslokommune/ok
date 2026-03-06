package workflow

import (
	"github.com/oslokommune/ok/pkg/workflow"
	"github.com/spf13/cobra"
)

var (
	appFlagAccountID           string
	appFlagRegion              string
	appFlagType                string
	appFlagDevEnvironmentName  string
	appFlagProdEnvironmentName string
	appFlagVarFiles            []string
)

// AppInitCommand initializes CI/CD workflow files for an application repository.
var AppInitCommand = &cobra.Command{
	Use:   "init <app-name>",
	Short: "Initialize CI/CD workflow files for an application repository",
	Long: `Initialize GitHub Actions workflow files for an application repository
using the golden-path-boilerplate app-cicd template.

This command runs boilerplate to download and render workflow templates.`,
	Example: `  # Basic initialization
  ok workflow app init my-app

  # For a repo that also contains infrastructure
  ok workflow app init my-app --type=app-with-iac

  # With AWS account and region
  ok workflow app init my-app --account-id 123456789012 --region eu-west-1`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := workflow.ValidateAppType(appFlagType); err != nil {
			return err
		}

		return workflow.RunAppInit(workflow.AppInitOptions{
			AppName:             args[0],
			AppType:             workflow.AppType(appFlagType),
			AccountID:           appFlagAccountID,
			Region:              appFlagRegion,
			DevEnvironmentName:  appFlagDevEnvironmentName,
			ProdEnvironmentName: appFlagProdEnvironmentName,
			VarFiles:            appFlagVarFiles,
		})
	},
}

func init() {
	AppInitCommand.Flags().StringVar(&appFlagAccountID, "account-id", "", "AWS account ID")
	AppInitCommand.Flags().StringVar(&appFlagRegion, "region", "", "AWS region")
	AppInitCommand.Flags().StringVar(&appFlagType, "type", "", "Repository type variant (valid: app-with-iac)")
	AppInitCommand.Flags().StringVar(&appFlagDevEnvironmentName, "dev-env-name", "", "Name of the dev environment, used in AWS resource names")
	AppInitCommand.Flags().StringVar(&appFlagProdEnvironmentName, "prod-env-name", "", "Name of the prod environment, used in AWS resource names")
	AppInitCommand.Flags().StringArrayVar(&appFlagVarFiles, "var-file", nil, "Path to a boilerplate variable file (repeatable)")
}
