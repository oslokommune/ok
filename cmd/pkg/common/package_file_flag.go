package common

import (
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/spf13/cobra"
)

func AddPackageFileFlag(cmd *cobra.Command, variable *string) {
	cmd.Flags().StringVarP(variable,
		common.FlagNamePackagesFile,
		"f",
		"packages.yml",
		"Set the path to the package manifest file",
	)
}
