package workflow

import (
	"github.com/oslokommune/ok/pkg/workflow"
	"github.com/spf13/cobra"
)

var (
	iacFlagDevAccountID        string
	iacFlagProdAccountID       string
	iacFlagDevRegion           string
	iacFlagProdRegion          string
	iacFlagDevEnvironmentName  string
	iacFlagProdEnvironmentName string
)

// IacInitCommand initializes CI/CD workflow files for a terraform-iac repository.
var IacInitCommand = &cobra.Command{
	Use:   "init",
	Short: "Initialize CI/CD workflow files for a terraform-iac repository",
	Long: `Initialize GitHub Actions workflow files for a pure infrastructure repository
using the golden-path-boilerplate terraform-iac template.

This command runs boilerplate to download and render workflow templates.`,
	Example: `  # Basic initialization
  ok workflow iac init

  # With AWS accounts and regions
  ok workflow iac init --dev-account-id 111111111111 --prod-account-id 222222222222 --dev-region eu-west-1 --prod-region eu-west-1

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return workflow.RunIacInit(workflow.IacInitOptions{
			DevAccountID:        iacFlagDevAccountID,
			ProdAccountID:       iacFlagProdAccountID,
			DevRegion:           iacFlagDevRegion,
			ProdRegion:          iacFlagProdRegion,
			DevEnvironmentName:  iacFlagDevEnvironmentName,
			ProdEnvironmentName: iacFlagProdEnvironmentName,
		})
	},
}

func init() {
	IacInitCommand.Flags().StringVar(&iacFlagDevAccountID, "dev-account-id", "", "AWS account ID for the dev environment")
	IacInitCommand.Flags().StringVar(&iacFlagProdAccountID, "prod-account-id", "", "AWS account ID for the prod environment")
	IacInitCommand.Flags().StringVar(&iacFlagDevRegion, "dev-region", "", "AWS region for the dev environment")
	IacInitCommand.Flags().StringVar(&iacFlagProdRegion, "prod-region", "", "AWS region for the prod environment")
	IacInitCommand.Flags().StringVar(&iacFlagDevEnvironmentName, "dev-env-name", "", "Name of the dev environment, used in AWS resource names")
	IacInitCommand.Flags().StringVar(&iacFlagProdEnvironmentName, "prod-env-name", "", "Name of the prod environment, used in AWS resource names")
}
