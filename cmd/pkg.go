package cmd

import (
	"github.com/oslokommune/ok/internal/pkg"
	"github.com/spf13/cobra"
)

var pkgCommand = &cobra.Command{
	Use:           "pkg",
	Short:         "Run pkg",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg.Get()
		err := pkg.Get()
		if err != nil {
			return err
		}
		return nil
	},
}
