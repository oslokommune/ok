package cmd

import (
	"github.com/oslokommune/ok/pkg/bootstrap"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(bootstrapCommand)
}

var bootstrapCommand = &cobra.Command{
	Use:   "bootstrap",
	Short: "This command will create the necessary S3 bucket and DynamoDB table that will be used to store Terraform state.",
	Run: func(cmd *cobra.Command, args []string) {
		//fullArgs := append([]string{"bootstrap"}, args...)
		//scriptrunner.RunScript("ok.sh", fullArgs)
		bootstrap.Entry()
	},
}
