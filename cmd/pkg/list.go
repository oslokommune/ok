package pkg

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/list"
	"github.com/spf13/cobra"
)

var ListCommand = &cobra.Command{
	Use:           "list",
	Short:         "List all defined packages.",
	Example:       `ok pkg list`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		packages, err := list.Run(PackagesManifestFilename)
		if err != nil {
			return err
		}

		for _, p := range packages {
			fmt.Println(p.String())
		}

		return nil
	},
}
