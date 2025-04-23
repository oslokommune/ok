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
		Run: func(cmd *cobra.Command, args []string) {
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

			// Print the boilerplate args that will be run
			fmt.Println("Boilerplate Command Arguments:")
			for _, config := range mergedConfigs {
				args := pk.BuildBoilerplateArgs(config)
				fmt.Printf("- %v\n", args)
			}
		},
	}

	return cmd
}
