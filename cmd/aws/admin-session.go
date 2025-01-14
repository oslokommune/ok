package aws

import (
	"os"
	"strconv"

	"github.com/oslokommune/ok/pkg/aws"
	"github.com/spf13/cobra"
)

var (
	startShell bool
	verbose    bool
)

var AdminSessionCommand = &cobra.Command{
	Use:   "admin-session",
	Short: "Start an admin session to an AWS account",

	RunE: func(cmd *cobra.Command, args []string) error {
		return aws.StartAdminSession(startShell, verbose)
	},
}

func init() {
	AdminSessionCommand.Flags().BoolVarP(&startShell, "start-shell", "s", false, "Start a working shell to execute AWS commands")
	AdminSessionCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Check for environment variable
	if envVerbose, exists := os.LookupEnv("OK_AWS_ADMIN_SESSION_VERBOSE"); exists {
		if parsedVerbose, err := strconv.ParseBool(envVerbose); err == nil {
			verbose = parsedVerbose
		}
	}

	// CLI flag takes precedence over environment variable
	AdminSessionCommand.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "Enable verbose output")
}

