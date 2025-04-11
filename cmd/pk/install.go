package pk

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Hello World install command",
		Long:  "This is a simple Hello World install command.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello World from the install command!")
		},
	}

	return cmd
}
