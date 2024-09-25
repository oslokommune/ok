package aws

import (
	"github.com/oslokommune/ok/pkg/aws"
	"github.com/spf13/cobra"
)

var startShell bool

var AdminSessionCommand = &cobra.Command{
	Use:   "admin-session",
	Short: "Start an admin session to an AWS account",

	RunE: func(cmd *cobra.Command, args []string) error {
		return aws.StartAdminSession(startShell)
	},
}

func init() {
	AdminSessionCommand.Flags().BoolVarP(&startShell, "start-shell", "s", false, "Start a working shell to execute AWS commands")
}
