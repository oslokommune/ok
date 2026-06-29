package pk

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/oslokommune/ok/pkg/pk"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	var configFile string
	var baseOutputFolder string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a .ok configuration directory",
		Long: `Creates a .ok directory with a config file.

If no flags are provided, prompts interactively for:
- Config file name (e.g., dev.yaml, prod.yaml, or config.yaml)
- Base output folder (e.g., stacks/dev, stacks/prod, or .)

The config file includes sensible defaults for the common section.`,
		Example: `  ok pk init                                    # interactive
  ok pk init --file dev.yaml --base stacks/dev  # non-interactive
  ok pk init --file prod.yaml --base stacks/prod`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			repoRoot, err := pk.RepoRoot(ctx)
			if err != nil {
				return errors.Join(fmt.Errorf("finding repository root"), err)
			}

			okDir := filepath.Join(repoRoot, ".ok")

			opts := pk.InitOptions{
				ConfigFileName:   configFile,
				BaseOutputFolder: baseOutputFolder,
			}

			// If neither flag provided, prompt interactively
			if configFile == "" && baseOutputFolder == "" {
				promptedOpts, err := promptInitOptions()
				if err != nil {
					return err
				}
				opts = promptedOpts
			}

			if err := pk.Init(okDir, opts); err != nil {
				return err
			}

			configPath := opts.ConfigFileName
			if configPath == "" {
				configPath = pk.DefaultConfigFileName
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\n✅ Created %s\n", filepath.Join(okDir, configPath))
			fmt.Fprintf(cmd.OutOrStdout(), "\nNext steps:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ok pk add <template> [subfolder]  - Add a template\n")
			fmt.Fprintf(cmd.OutOrStdout(), "  ok pk install                     - Install templates\n")

			return nil
		},
	}

	cmd.Flags().StringVar(&configFile, "file", "", "config file name (e.g., dev.yaml)")
	cmd.Flags().StringVar(&baseOutputFolder, "base", "", "base output folder (e.g., stacks/dev)")

	return cmd
}

func promptInitOptions() (pk.InitOptions, error) {
	var configFileName string
	var baseOutputFolder string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Config file name").
				Description("Name for the config file (e.g., dev.yaml, prod.yaml)").
				Placeholder("config.yaml").
				Value(&configFileName),

			huh.NewInput().
				Title("Base output folder").
				Description("Where templates will be generated (e.g., stacks/dev)").
				Placeholder(".").
				Value(&baseOutputFolder),
		),
	)

	err := form.Run()
	if err != nil {
		return pk.InitOptions{}, fmt.Errorf("prompting for options: %w", err)
	}

	// Clean up inputs
	configFileName = strings.TrimSpace(configFileName)
	baseOutputFolder = strings.TrimSpace(baseOutputFolder)

	// Ensure .yaml extension
	if configFileName != "" && !strings.HasSuffix(configFileName, ".yaml") && !strings.HasSuffix(configFileName, ".yml") {
		configFileName = configFileName + ".yaml"
	}

	return pk.InitOptions{
		ConfigFileName:   configFileName,
		BaseOutputFolder: baseOutputFolder,
	}, nil
}
