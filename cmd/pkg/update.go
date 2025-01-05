package pkg

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/config"
	"github.com/oslokommune/ok/pkg/pkg/install"
	"github.com/oslokommune/ok/pkg/pkg/install/interactive"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/update"
	"github.com/spf13/cobra"
)

type GitHubReleases interface {
	GetLatestReleases() (map[string]string, error)
}

func NewUpdateCommand(ghReleases GitHubReleases, downloader config.FileDownloader) *cobra.Command {
	var flagDisableManifestUpdate bool
	var flagUpdateCommandUpdateSchema bool
	var flagMigrateConfig bool

	updater := update.NewUpdater(ghReleases, downloader)

	cmd := &cobra.Command{
		Use:   "update [outputFolder ...]",
		Short: "Update Boilerplate package manifest and package configuration files",
		Long: `Update Boilerplate package manifest and package configuration files.

` + InstallUpdateArgumentDescription,
		Example: `ok pkg update
ok pkg update my-package
`,
		ValidArgsFunction: updateTabCompletion,
		RunE: func(cmd *cobra.Command, outputFolders []string) error {
			if flagInteractive && len(outputFolders) > 0 {
				return fmt.Errorf("cannot use both --interactive and outputFolder arguments")
			}

			var packages []common.Package
			var err error

			manifest, err := common.LoadPackageManifest(common.PackagesManifestFilename)
			if err != nil {
				return fmt.Errorf("loading package manifest: %w", err)
			}

			// Select packages
			switch {
			case len(outputFolders) > 0:
				// Use output folders to determine which packages to install
				packages = install.FindPackagesFromOutputFolders(manifest.Packages, outputFolders)

			case flagInteractive:
				// Use interactive mode to determine which packages to install
				packages, err = interactive.SelectPackages(manifest, "update")
				if err != nil {
					return fmt.Errorf("selecting packages: %w", err)
				}

				if len(packages) == 0 {
					fmt.Println("No packages selected. Remember to use space (or x) to select package(s) to install.")
					return nil
				}

			default:
				packages = manifest.Packages
			}

			opts := update.Options{
				DisableManifestUpdate: flagDisableManifestUpdate,
				MigrateConfig:         flagMigrateConfig,
				UpdateSchemaConfig:    flagUpdateCommandUpdateSchema,
			}

			err = updater.Run(common.PackagesManifestFilename, packages, opts)
			if err != nil {
				return fmt.Errorf("updating packages: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&flagInteractive,
		FlagInteractiveName, FlagInteractiveShorthand, false, FlagInteractiveUsage)

	cmd.Flags().BoolVar(&flagDisableManifestUpdate,
		"disable-manifest-update",
		false,
		"Disable package manifest version updates (useful when using an external dependency manager like Renovate)",
	)

	cmd.Flags().BoolVar(&flagUpdateCommandUpdateSchema,
		"update-schema",
		true,
		"Update the JSON schema for affected packages")

	cmd.Flags().BoolVar(&flagMigrateConfig,
		"migrate-config",
		true,
		"Automatically migrate package configuration files to the latest version, if possible")

	return cmd
}

func updateTabCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// No completion if there are more than one argument
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	manifest, err := common.LoadPackageManifest(common.PackagesManifestFilename)
	if err != nil {
		cmd.PrintErrf("failed to load package manifest: %s\n", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, p := range manifest.Packages {
		if strings.HasPrefix(p.OutputFolder, toComplete) {
			completions = append(completions, p.OutputFolder)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
