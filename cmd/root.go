package cmd

import (
	"github.com/oslokommune/ok/scriptrunner"
	"github.com/oslokommune/ok/toggle"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "ok",
		Short: "Command Runner is a simple tool to run a script with subcommands",
	}

	rootCmd.AddCommand(newBootstrapCommand(),
		newScaffoldCommand(),
		newEnvCommand(),
		newEnvarsCommand(),
		newGetTemplateCommand(),
		newForwardCommand(),
		newVersionCommand(),
		newAssumeCommand())

	return rootCmd
}

func newBootstrapCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap",
		Short: "This command will create the necessary S3 bucket and DynamoDB table that will be used to store Terraform state.",
		Run: func(cmd *cobra.Command, args []string) {
			fullArgs := append([]string{"scaffold"}, args...)
			scriptrunner.RunScript("ok.sh", fullArgs)
		},
	}
}

func newScaffoldCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "scaffold",
		Short: "Creates a new Terraform project with a _config.tf, _variables.tf, _versions.tf and _config.auto.tfvars.json file based on values configured in env.yml.",
		Run: func(cmd *cobra.Command, args []string) {
			fullArgs := append([]string{"scaffold"}, args...)
			scriptrunner.RunScript("ok.sh", fullArgs)
		},
	}
}

func newEnvCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "Creates a new env.yml file with placeholder values.",
		Run: func(cmd *cobra.Command, args []string) {
			fullArgs := append([]string{"env"}, args...)
			scriptrunner.RunScript("ok.sh", fullArgs)
		},
	}
}

func newEnvarsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "envars",
		Short: "Exports the values in env.yml as environment variables.",
		Run: func(cmd *cobra.Command, args []string) {
			fullArgs := append([]string{"envars"}, args...)
			scriptrunner.RunScript("ok.sh", fullArgs)
		},
	}
}

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

func newForwardCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "forward",
		Short: "Starts a port forwarding session to a database.",
		Run: func(cmd *cobra.Command, args []string) {
			scriptrunner.RunScript("port-forward.sh", args)
		},
	}
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of the ok tool and the current latest version available.",
		Run: func(cmd *cobra.Command, args []string) {
			fullArgs := append([]string{"version"}, args...)
			scriptrunner.RunScript("ok.sh", fullArgs)
		},
	}
}

func newAssumeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "assume",
		Short: "Toggle assume_cd_role in app stack",
		Run: func(cmd *cobra.Command, args []string) {
			toggle.Assume()
		},
	}
}
