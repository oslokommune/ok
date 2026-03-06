package workflow

import "github.com/spf13/cobra"

// AppCommand is the parent command group for app workflow commands.
var AppCommand = &cobra.Command{
	Use:   "app",
	Short: "Application workflow commands.",
}
