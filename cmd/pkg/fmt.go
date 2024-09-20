package pkg

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/format"
	"github.com/spf13/cobra"
)

func init() {
}

var FmtCommand = &cobra.Command{
	Use:           "fmt",
	Short:         "Format the package manifest file.",
	Example:       `ok pkg fmt`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := format.Run(common.PackagesManifestFilename)
		if err != nil {
			return fmt.Errorf("formatting package manifest file packages: %w", err)
		}

		return nil
	},
}
