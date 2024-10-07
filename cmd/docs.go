package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCommand = &cobra.Command{
	Use:                   "docs",
	Short:                 "Generates ok's command line docs.",
	Long:                  "Generates Markdown documentation for all commands in the CLI.",
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	Hidden:                true,
	Args:                  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		rootCmd.Root().DisableAutoGenTag = true
		return doc.GenMarkdownTree(rootCmd.Root(), "docs")
	},
}

func init() {
	rootCmd.AddCommand(docsCommand)
}
