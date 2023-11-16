package cmd

import (
	"github.com/oslokommune/ok/scriptrunner"
	"github.com/spf13/cobra"
)

func newGetTemplateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get-template",
		Short: "Downloads a template from the golden-path-iac repository.",
		Run: func(cmd *cobra.Command, args []string) {
			fullArgs := append([]string{"get-template"}, args...)
			scriptrunner.RunScript("ok.sh", fullArgs)
		},
	}
}
