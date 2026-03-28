package pk

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/oslokommune/ok/pkg/pk"
	"github.com/spf13/cobra"
)

func NewAddCommand(ghReleases pk.GitHubReleases) *cobra.Command {
	var ref string
	var file string

	cmd := &cobra.Command{
		Use:   "add [template] [subfolder]",
		Short: "Add a template to the .ok configuration",
		Long: `Add a template to the .ok configuration file.

If no template is specified, prompts interactively to select from available templates.
The template version is fetched from the latest GitHub release unless --ref is specified.
If no subfolder is provided, prompts for one (defaults to template name).`,
		Example: `  ok pk add                         # interactive
  ok pk add app                     # prompts for subfolder
  ok pk add app my-app              # non-interactive
  ok pk add app my-app --ref v10.0.0
  ok pk add networking --file .ok/infra.yaml`,
		Args: cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var templateName string
			var subfolder string

			if len(args) >= 1 {
				templateName = args[0]
			}
			if len(args) >= 2 {
				subfolder = args[1]
			}

			okDir, err := pk.OkDir(ctx)
			if err != nil {
				return errors.Join(fmt.Errorf("locating .ok directory"), err)
			}

			// Check if .ok directory exists, offer to initialize if not
			exists, err := dirExists(okDir)
			if err != nil {
				return errors.Join(fmt.Errorf("checking .ok directory"), err)
			}

			if !exists {
				shouldInit, err := promptInit()
				if err != nil {
					return err
				}
				if !shouldInit {
					return fmt.Errorf("cannot add template without .ok directory")
				}

				// Prompt for init options
				initOpts, err := promptInitOptions()
				if err != nil {
					return err
				}

				if err := pk.Init(okDir, initOpts); err != nil {
					return errors.Join(fmt.Errorf("initializing .ok directory"), err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "✅ Initialized .ok directory\n\n")

				// If a specific config file was created, use it
				if initOpts.ConfigFileName != "" {
					file = filepath.Join(okDir, initOpts.ConfigFileName)
				}
			}

			// If no template specified, prompt to select one
			if templateName == "" {
				selected, err := promptTemplateSelection(ghReleases)
				if err != nil {
					return err
				}
				templateName = selected
			}

			// If no subfolder specified, prompt for one
			if subfolder == "" {
				prompted, err := promptSubfolder(templateName)
				if err != nil {
					return err
				}
				subfolder = prompted
			}

			opts := pk.AddOptions{
				TemplateName: templateName,
				Subfolder:    subfolder,
				Ref:          ref,
				ConfigFile:   file,
			}

			if err := pk.Add(okDir, opts, ghReleases); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&ref, "ref", "", "specific version (e.g., v10.0.0)")
	cmd.Flags().StringVar(&file, "file", "", "target config file (e.g., .ok/config.yaml)")

	return cmd
}

func promptInit() (bool, error) {
	var confirm bool

	err := huh.NewConfirm().
		Title("No .ok directory found. Initialize now?").
		Affirmative("Yes").
		Negative("No").
		Value(&confirm).
		Run()

	if err != nil {
		return false, fmt.Errorf("prompting for init: %w", err)
	}

	return confirm, nil
}

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func promptTemplateSelection(ghReleases pk.GitHubReleases) (string, error) {
	fmt.Println("Fetching available templates...")

	releases, err := ghReleases.GetLatestReleases()
	if err != nil {
		return "", fmt.Errorf("fetching releases: %w", err)
	}

	// Build options from available templates
	options := make([]huh.Option[string], 0, len(releases))
	for template, version := range releases {
		label := fmt.Sprintf("%s (%s)", template, version)
		options = append(options, huh.NewOption(label, template))
	}

	var selected string
	err = huh.NewSelect[string]().
		Title("Select a template").
		Options(options...).
		Value(&selected).
		Run()

	if err != nil {
		return "", fmt.Errorf("selecting template: %w", err)
	}

	return selected, nil
}

func promptSubfolder(templateName string) (string, error) {
	var subfolder string

	err := huh.NewInput().
		Title("Subfolder name").
		Description("Where the template will be generated").
		Placeholder(templateName).
		Value(&subfolder).
		Run()

	if err != nil {
		return "", fmt.Errorf("prompting for subfolder: %w", err)
	}

	// Use template name as default if empty
	if subfolder == "" {
		subfolder = templateName
	}

	return subfolder, nil
}
