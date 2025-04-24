package pk

import (
	"fmt"

	"github.com/oslokommune/ok/pkg/pk"
	"github.com/spf13/cobra"
)

func NewInstallCommand() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Generate project boilerplate",
		Long:  "Reads .ok configuration files, merges them and runs the boilerplate generator.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			okDir, err := pk.OkDir(ctx)
			if err != nil {
				return fmt.Errorf("locating .ok directory: %w", err)
			}

			repoDir, err := pk.RepoRoot(ctx)
			if err != nil {
				return fmt.Errorf("finding repository root: %w", err)
			}

			configs, err := pk.LoadConfigs(okDir)
			if err != nil {
				return fmt.Errorf("loading configs: %w", err)
			}

			effectiveCfgs, err := pk.ApplyCommon(configs)
			if err != nil {
				return fmt.Errorf("applying common configs: %w", err)
			}

			for _, cfg := range effectiveCfgs {
				args, err := pk.BuildBoilerplateArgs(cfg)
				if err != nil {
					return fmt.Errorf("building boilerplate args: %w", err)
				}

				if dryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "dry-run: %v\n", args)
					continue
				}

				if err := pk.RunBoilerplateCommand(ctx, args, repoDir); err != nil {
					return fmt.Errorf("running boilerplate command: %w", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false,
		"print the generated boilerplate command without executing it")

	return cmd
}
