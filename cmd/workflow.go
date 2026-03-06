package cmd

import "github.com/spf13/cobra"

var workflowCommand = &cobra.Command{
	Use:   "workflow",
	Short: "Initialize CI/CD workflow files.",
}
