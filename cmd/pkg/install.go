package pkg

import (
	"github.com/oslokommune/ok/pkg/pkg/install"
	"github.com/spf13/cobra"
)

var InstallCommand = &cobra.Command{
	Use:           "install",
	Short:         "Run install",
	SilenceErrors: true, // TODO true
	RunE: func(cmd *cobra.Command, args []string) error {
		err := install.Run()
		if err != nil {
			return err
		}

		return nil
	},
}
