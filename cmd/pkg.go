package cmd

import (
	"github.com/oslokommune/ok/internal/pkg"
	"github.com/spf13/cobra"
)

var pkgCommand = &cobra.Command{
	Use:   "pkg",
	Short: "Run pkg",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Get()
	},
}
