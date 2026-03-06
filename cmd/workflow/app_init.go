package workflow

import (
	"github.com/oslokommune/ok/pkg/workflow"
	"github.com/spf13/cobra"
)

var (
	appFlagDevAccountID        string
	appFlagProdAccountID       string
	appFlagDevRegion           string
	appFlagProdRegion          string
	appFlagType                string
	appFlagDevEnvironmentName  string
	appFlagProdEnvironmentName string
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

  # With AWS accounts and regions
  ok workflow app init my-app --dev-account-id 111111111111 --prod-account-id 222222222222 --dev-region eu-west-1 --prod-region eu-west-1`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := workflow.ValidateAppType(appFlagType); err != nil {
			return err
		}

		return workflow.RunAppInit(workflow.AppInitOptions{
			AppName:             args[0],
			AppType:             workflow.AppType(appFlagType),
			DevAccountID:        appFlagDevAccountID,
			ProdAccountID:       appFlagProdAccountID,
			DevRegion:           appFlagDevRegion,
			ProdRegion:          appFlagProdRegion,
			DevEnvironmentName:  appFlagDevEnvironmentName,
			ProdEnvironmentName: appFlagProdEnvironmentName,
		})
	},
}

func init() {
	AppInitCommand.Flags().StringVar(&appFlagDevAccountID, "dev-account-id", "", "AWS account ID for the dev environment")
	AppInitCommand.Flags().StringVar(&appFlagProdAccountID, "prod-account-id", "", "AWS account ID for the prod environment")
	AppInitCommand.Flags().StringVar(&appFlagDevRegion, "dev-region", "", "AWS region for the dev environment")
	AppInitCommand.Flags().StringVar(&appFlagProdRegion, "prod-region", "", "AWS region for the prod environment")
	AppInitCommand.Flags().StringVar(&appFlagType, "type", "", "Repository type variant (valid: app-with-iac)")
	AppInitCommand.Flags().StringVar(&appFlagDevEnvironmentName, "dev-env-name", "", "Name of the dev environment, used in AWS resource names")
	AppInitCommand.Flags().StringVar(&appFlagProdEnvironmentName, "prod-env-name", "", "Name of the prod environment, used in AWS resource names")
}
