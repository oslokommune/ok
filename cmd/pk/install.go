package pk

import (
	"errors"
	"fmt"
	"os"

	"github.com/oslokommune/ok/pkg/pk"
	"github.com/spf13/cobra"
)

func NewInstallCommand() *cobra.Command {
	var dryRun bool
	var all bool

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Generate project boilerplate",
		Long: `Reads .ok configuration files, merges them and runs the boilerplate generator.

By default, the command is context-aware:
- If run from a subfolder that matches a template's output path, only that template is installed
- If run from the repository root (or no match), an interactive picker is shown

Use --all to bypass filtering and install all templates.`,
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

			// Determine which templates to install
			templatesToInstall := effectiveCfgs

			if !all {
				cwd, err := os.Getwd()
				if err != nil {
					return errors.Join(fmt.Errorf("getting current working directory"), err)
				}

				// Try to filter by current directory
				matched := pk.FilterTemplatesByWorkingDir(effectiveCfgs, cwd, repoDir)

				if len(matched) > 0 {
					// Found matching templates based on cwd
					templatesToInstall = matched
				} else {
					// No match (e.g., at repo root) - show interactive picker
					selected, err := pk.SelectTemplatesInteractively(effectiveCfgs)
					if err != nil {
						return errors.Join(fmt.Errorf("selecting templates"), err)
					}
					if len(selected) == 0 {
						fmt.Fprintln(cmd.OutOrStdout(), "No templates selected.")
						return nil
					}
					templatesToInstall = selected
				}
			}

			for _, cfg := range templatesToInstall {
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
	cmd.Flags().BoolVar(&all, "all", false,
		"install all templates without filtering or prompting")

	return cmd
}
