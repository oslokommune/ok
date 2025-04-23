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

			// Print useful information about the merged configurations
			fmt.Println("Merged Template Configurations:")
			for _, config := range mergedConfigs {
				fmt.Printf("- Name: %s, Repo: %s, Ref: %s, Path: %s, Subfolder: %s, VarFiles: %v\n",
					config.Name, config.Repo, config.Ref, config.Path, config.Subfolder, config.VarFiles)
			}
		},
	}

	return cmd
}
