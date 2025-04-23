package pk

import (
	"fmt"

	"github.com/oslokommune/ok/pkg/pk"
	"github.com/spf13/cobra"
)

func NewInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Hello World install command",
		Long:  "This is a simple Hello World install command.",
	}

	cmd.Flags().Bool("dry-run", false, "Print the boilerplate command arguments without executing them")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		okDir, err := pk.OkDir()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Load configs from the .ok directory
		configs, err := pk.LoadConfigs(okDir)
		if err != nil {
			fmt.Printf("Error loading configs: %v\n", err)
			return
		}

		// Generate merged template configurations
		mergedConfigs, err := pk.ApplyCommon(configs)
		if err != nil {
			fmt.Printf("Error generating merged configs: %v\n", err)
			return
		}

		// Check if dry-run flag is set
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		for _, config := range mergedConfigs {
			args := pk.BuildBoilerplateArgs(config)
			if dryRun {
				fmt.Printf("Dry-run: %v\n", args)
			} else {
				fmt.Printf("Running boilerplate command with args: %v\n", args)
				err := pk.RunBoilerplateCommand(args, okDir)
				if err != nil {
					fmt.Printf("Error running boilerplate command: %v\n", err)
					return
				}
			}
		}
	}

	return cmd
}
