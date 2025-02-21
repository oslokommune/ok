package aws

import (
	"github.com/oslokommune/ok/pkg/aws/config"
	"github.com/spf13/cobra"
)

var (
	ssoStarturl string
	ssoRegion   string
	region      string
	sessionName string
	template    string
)

var ConfigGeneratorCommand = &cobra.Command{
	Use:   "generate",
	Short: "Generate AWS CLI configuration for AWS IAM Identity Center roles",
	Long: `Generate AWS CLI configuration for AWS IAM Identity Center roles.

Profile name template supports the following variables:
  {{.SessionName}}  - The name of the SSO session
  {{.AccountName}}  - The name of the AWS account
  {{.AccountID}}    - The ID of the AWS account
  {{.RoleName}}     - The name of the IAM role

Example:
  ok aws generate \
    --sso-start-url "https://my-sso.awsapps.com/start" \
    --sso-region "eu-west-1" \
    --template "ok-{{.AccountName}}-{{.RoleName}}"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if region == "" {
			region = ssoRegion
		}

		return config.Generate(config.Options{
			SsoStartUrl: ssoStarturl,
			SsoRegion:   ssoRegion,
			Region:      region,
			SessionName: sessionName,
			Template:    template,
		})
	},
}

func init() {
	ConfigGeneratorCommand.Flags().StringVar(&ssoStarturl, "sso-start-url", "", "The start URL of the AWS IAM Identity Center instance")
	ConfigGeneratorCommand.Flags().StringVar(&ssoRegion, "sso-region", "", "The region of the AWS IAM Identity Center instance")
	ConfigGeneratorCommand.Flags().StringVar(&region, "region", "", "The default region for generated profiles (defaults to sso-region if not set)")
	ConfigGeneratorCommand.Flags().StringVar(&sessionName, "session-name", "", "An optional name of the SSO session (defaults to the AWS IAM Identity Center instance identifier if not set)")
	ConfigGeneratorCommand.Flags().StringVar(&template, "template", "", "Go string template for generating profile names")
	ConfigGeneratorCommand.MarkFlagRequired("sso-start-url")
	ConfigGeneratorCommand.MarkFlagRequired("sso-region")
}
