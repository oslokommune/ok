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
			okDir, err := pk.GetOkDirPath()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("The .ok directory is located at: %s\n", okDir)
		},
	}

	return cmd
}
