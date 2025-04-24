package pk

import (
	"errors"
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
				return errors.Join(fmt.Errorf("locating .ok directory"), err)
			}

			repoDir, err := pk.RepoRoot(ctx)
			if err != nil {
				return errors.Join(fmt.Errorf("finding repository root"), err)
			}

			configs, err := pk.LoadConfigs(okDir)
			if err != nil {
				return errors.Join(fmt.Errorf("loading configs"), err)
			}

			effectiveCfgs, err := pk.ApplyCommon(configs)
			if err != nil {
				return errors.Join(fmt.Errorf("applying common configs"), err)
			}

			for _, cfg := range effectiveCfgs {
				args := pk.BuildBoilerplateArgs(cfg)

				if dryRun {
					fmt.Fprintf(cmd.OutOrStdout(), "dry-run: %v\n", args)
					continue
				}

				if err := pk.RunBoilerplateCommand(ctx, args, repoDir); err != nil {
					return errors.Join(fmt.Errorf("running boilerplate command"), err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false,
		"print the generated boilerplate command without executing it")

	return cmd
}
