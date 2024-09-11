package cmd

import (
	"github.com/oslokommune/ok/internal/scriptrunner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getTemplateCommand)
}

var getTemplateCommand = &cobra.Command{
	Use:        "get-template",
	Short:      "Downloads a template from the `golden-path-iac` repository.",
	Deprecated: "and will be removed in a future release. Use the `pkg` command instead.",
	Run: func(cmd *cobra.Command, args []string) {
		fullArgs := append([]string{"get-template"}, args...)
		scriptrunner.RunScript("ok.sh", fullArgs)
	},
}
