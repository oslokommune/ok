package workflow

import "github.com/spf13/cobra"

// IacCommand is the parent command group for IAC workflow commands.
var IacCommand = &cobra.Command{
	Use:   "iac",
	Short: "Infrastructure-as-code workflow commands.",
}
