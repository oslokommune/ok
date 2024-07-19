package pkg

import (
	"github.com/oslokommune/ok/pkg/pkg/update/select_pkg"
	"github.com/spf13/cobra"
)

var UpdateCommand = &cobra.Command{
	Use:           "update",
	Short:         "Update Boilerplate package manifest.",
	Example:       `ok pkg update`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := select_pkg.Run()
		if err != nil {
			return err
		}

		return nil
	},
}
