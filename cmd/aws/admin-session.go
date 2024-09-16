package aws

import (
	"github.com/oslokommune/ok/pkg/aws"
	"github.com/spf13/cobra"
)

var AdminSessionCommand = &cobra.Command{
	Use:   "admin-session",
	Short: "Start an admin session to an AWS account",

	RunE: func(cmd *cobra.Command, args []string) error {
		return aws.StartAdminSession()
	},
}
