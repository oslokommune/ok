package pkg

import (
	"github.com/oslokommune/ok/pkg/pkg/install"
	"github.com/spf13/cobra"
)

var InstallCommand = &cobra.Command{
	Use:           "install [stackName1, stackName2, ...]",
	Short:         "Install or update Boilerplate packages",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := install.Run(args)
		if err != nil {
			return err
		}

		return nil
	},
}
