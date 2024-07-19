package pkg

import (
	"github.com/oslokommune/ok/pkg/pkg/update"
	"github.com/spf13/cobra"
)

var UpdateCommand = &cobra.Command{
	Use:           "update",
	Short:         "Update Boilerplate package manifest.",
	Example:       `ok pkg update`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := update.Run(PackagesManifestFilename)
		if err != nil {
			return err
		}

		return nil
	},
}
