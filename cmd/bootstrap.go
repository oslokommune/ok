package cmd

import (
	"github.com/oslokommune/ok/internal/scriptrunner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(bootstrapCommand)
}

var bootstrapCommand = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap code for an S3 bucket and DynamoDB table to store Terraform state.",
	Run: func(cmd *cobra.Command, args []string) {
		fullArgs := append([]string{"bootstrap"}, args...)
		scriptrunner.RunScript("ok.sh", fullArgs)
	},
}
