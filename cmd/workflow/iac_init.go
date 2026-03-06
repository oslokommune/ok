package workflow

import (
	"github.com/oslokommune/ok/pkg/workflow"
	"github.com/spf13/cobra"
)

var (
	iacFlagAccountID           string
	iacFlagRegion              string
	iacFlagDevEnvironmentName  string
	iacFlagProdEnvironmentName string
	iacFlagVarFiles            []string
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

  # With AWS account and region
  ok workflow iac init --account-id 123456789012 --region eu-west-1

  # With a boilerplate variable file
  ok workflow iac init --var-file common-config.yml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return workflow.RunIacInit(workflow.IacInitOptions{
			AccountID:           iacFlagAccountID,
			Region:              iacFlagRegion,
			DevEnvironmentName:  iacFlagDevEnvironmentName,
			ProdEnvironmentName: iacFlagProdEnvironmentName,
			VarFiles:            iacFlagVarFiles,
		})
	},
}

func init() {
	IacInitCommand.Flags().StringVar(&iacFlagAccountID, "account-id", "", "AWS account ID")
	IacInitCommand.Flags().StringVar(&iacFlagRegion, "region", "", "AWS region")
	IacInitCommand.Flags().StringVar(&iacFlagDevEnvironmentName, "dev-env-name", "", "Name of the dev environment, used in AWS resource names")
	IacInitCommand.Flags().StringVar(&iacFlagProdEnvironmentName, "prod-env-name", "", "Name of the prod environment, used in AWS resource names")
	IacInitCommand.Flags().StringArrayVar(&iacFlagVarFiles, "var-file", nil, "Path to a boilerplate variable file (repeatable)")
}
