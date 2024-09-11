package cmd

import (
	"github.com/spf13/cobra"
)

var pkgCommand = &cobra.Command{
	Use:   "pkg",
	Short: "Group of package related commands for managing Boilerplate packages.",
}
