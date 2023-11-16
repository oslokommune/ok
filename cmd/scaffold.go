package cmd

import (
	"github.com/oslokommune/ok/scriptrunner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scaffoldCommand)
}

var scaffoldCommand = &cobra.Command{
	Use:   "scaffold",
	Short: "Creates a new Terraform project with a _config.tf, _variables.tf, _versions.tf and _config.auto.tfvars.json file based on values configured in env.yml.",
	Run: func(cmd *cobra.Command, args []string) {
		fullArgs := append([]string{"scaffold"}, args...)
		scriptrunner.RunScript("ok.sh", fullArgs)
	},
}
