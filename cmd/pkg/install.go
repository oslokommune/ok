package pkg

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/install/interactive"

	"github.com/oslokommune/ok/pkg/pkg/install"
	"github.com/spf13/cobra"
)

func NewInstallCommand() *cobra.Command {
	var flagRecursive bool

	cmd := &cobra.Command{
		Use:   "install [outputFolder ...]",
		Short: "Install or update Boilerplate packages.",
		Long: `Install or update Boilerplate packages.

` + InstallUpdateArgumentDescription,
		Example: `ok pkg install networking
ok pkg install networking my-app
BASE_URL=../boilerplate/terraform ok pkg install networking my-app
`,
		ValidArgsFunction: installTabCompletion,
		SilenceErrors:     true,
		RunE: func(cmd *cobra.Command, outputFolders []string) error {
			if len(outputFolders) > 0 && flagInteractive {
				return fmt.Errorf("cannot use both outputFolder arguments and --interactive arguments")
			}

			if len(outputFolders) > 0 && flagRecursive {
				return fmt.Errorf("cannot use both outputFolder arguments and --recursive arguments")
			}

			if flagInteractive && flagRecursive {
				return fmt.Errorf("cannot use both --interactive and --recursive arguments")
			}

			if flagRecursive {
				return runRecursiveInSubdirs(installWithBanner)
			} else {
				return installFromManifest(common.PackagesManifestFilename, outputFolders, ".")
			}
		},
	}

	cmd.Flags().BoolVarP(&flagInteractive,
		FlagInteractiveName, FlagInteractiveShorthand, false, FlagInteractiveUsage)

	//cmd.Flags().BoolVarP(&flagRecursive,
	//	"recursive",
	//	"r",
	//	false,
	//	"Install packages from manifests found in all subdirectories, but excluding the current directory.",
	//)

	addRecursiveFlagToCmd(cmd, &flagRecursive, "Install")

	return cmd
}

func installWithBanner(manifestPath string, manifestDir string, style lipgloss.Style) error {
	fmt.Println()
	fmt.Println(style.Render(fmt.Sprintf("Installing package manifest: %s", manifestPath)))
	fmt.Println()

	err := installFromManifest(manifestPath, []string{}, manifestDir)
	if err != nil {
		return fmt.Errorf("installing manifest %s: %w", manifestPath, err)
	}

	return nil
}

func installFromManifest(manifestFile string, outputFolders []string, workingDirectory string) error {
	var packages []common.Package
	var err error

	manifest, err := common.LoadPackageManifest(manifestFile)
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

	err = install.Run(packages, manifest, workingDirectory)
	if err != nil {
		return fmt.Errorf("installing packages: %w", err)
	}

	return nil
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
