package pkg

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/install"
	"github.com/oslokommune/ok/pkg/pkg/install/interactive"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/update"
	"github.com/spf13/cobra"
)

var flagUpdateCommandUpdateSchema bool

func init() {
	UpdateCommand.Flags().BoolVarP(&flagInteractive,
		FlagInteractiveName, FlagInteractiveShorthand, false, FlagInteractiveUsage)
	UpdateCommand.Flags().BoolVar(&flagUpdateCommandUpdateSchema,
		"update-schema", true, "Update the JSON schema for affected packages")
}

var UpdateCommand = &cobra.Command{
	Use:   "update [package-name]",
	Short: "Update Boilerplate package manifest",
	Long: `Update Boilerplate package manifest.
If a package name is provided, only that package will be updated.
If no package name is provided, all packages will be updated.`,
	Example: `  ok pkg update
  ok pkg update my-package`,
	Args:              cobra.MaximumNArgs(1),
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

		err = update.Run(common.PackagesManifestFilename, packages, flagUpdateCommandUpdateSchema)
		if err != nil {
			return fmt.Errorf("updating packages: %w", err)
		}

		return nil
	},
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
