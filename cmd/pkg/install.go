package pkg

import (
	"fmt"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/install/interactive"

	"github.com/oslokommune/ok/pkg/pkg/install"
	"github.com/spf13/cobra"
)

func init() {
	InstallCommand.Flags().BoolVarP(&flagInteractive,
		FlagInteractiveName, FlagInteractiveShorthand, false, FlagInteractiveUsage)
}

var InstallCommand = &cobra.Command{
	Use:   "install [outputFolder ...]",
	Short: "Install or update Boilerplate packages.",
	Long: `Install or update Boilerplate packages.

If no arguments are used, the command installs all the packages specified in the package manifest file.

If one or more output folders are specified, the command installs only the packages whose OutputFolder matches the specified folders. (OutputFolder is a field in the package manifest file.)

Set the environment variable BASE_URL to specify where package templates are downloaded from.
`,
	Example: `ok pkg install networking
ok pkg install networking my-app
BASE_URL=../boilerplate/terraform ok pkg install networking my-app
`,
	ValidArgsFunction: installTabCompletion,
	SilenceErrors:     true,
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
			packages, err = interactive.SelectPackages(manifest, "install")
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

		err = install.Run(packages, manifest)
		if err != nil {
			return fmt.Errorf("installing packages: %w", err)
		}

		return nil
	},
}

func installTabCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	manifest, err := common.LoadPackageManifest(common.PackagesManifestFilename)
	if err != nil {
		cmd.PrintErrf("failed to load package manifest: %s\n", err)
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, p := range manifest.Packages {
		if argsContainsElement(args, p.OutputFolder) {
			continue
		}
		if strings.HasPrefix(p.OutputFolder, toComplete) {
			completions = append(completions, p.OutputFolder)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
