package aws

import (
	"os"
	"strconv"

	"github.com/oslokommune/ok/pkg/aws"
	"github.com/spf13/cobra"
)

var (
	startShell bool
	verbosity  int
)

var AdminSessionCommand = &cobra.Command{
	Use:   "admin-session",
	Short: "Start an admin session to an AWS account",

	RunE: func(cmd *cobra.Command, args []string) error {
		return aws.StartAdminSession(startShell, verbosity)
	},
}

func init() {
	AdminSessionCommand.Flags().BoolVarP(&startShell, "start-shell", "s", false, "Start a working shell to execute AWS commands")
	AdminSessionCommand.Flags().IntVarP(&verbosity, "verbosity", "v", 1, "Set verbosity level (0-1)")

	// Check for environment variable
	if envVerbosity, exists := os.LookupEnv("OK_AWS_ADMIN_SESSION_VERBOSITY"); exists {
		if parsedVerbosity, err := strconv.Atoi(envVerbosity); err == nil {
			if parsedVerbosity >= 0 && parsedVerbosity <= 1 {
				verbosity = parsedVerbosity
			}
		}
	}

	// CLI flag takes precedence over environment variable
	AdminSessionCommand.Flags().IntVarP(&verbosity, "verbosity", "v", verbosity, "Set verbosity level (0-1)")
}

