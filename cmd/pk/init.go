package pk

import (
	"errors"
	"fmt"

	"github.com/oslokommune/ok/pkg/pk"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a .ok configuration directory",
		Long: `Creates a .ok directory with a default config.yaml file.

The config file includes sensible defaults for the common section:
- repo: git@github.com:oslokommune/golden-path-boilerplate.git
- non_interactive: true
- base_output_folder: "."`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			repoRoot, err := pk.RepoRoot(ctx)
			if err != nil {
				return errors.Join(fmt.Errorf("finding repository root"), err)
			}

			okDir := repoRoot + "/.ok"

			if err := pk.Init(okDir); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "✅ Initialized .ok directory at %s\n", okDir)
			fmt.Fprintf(cmd.OutOrStdout(), "\nNext steps:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ok pk add <template> [subfolder]  - Add a template\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ok pk install                     - Install templates\n")

			return nil
		},
	}

	return cmd
}
