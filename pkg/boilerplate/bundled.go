package boilerplate

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gruntwork-io/boilerplate/options"
	"github.com/gruntwork-io/boilerplate/templates"
	"github.com/gruntwork-io/boilerplate/variables"
	"github.com/oslokommune/ok/pkg/pkg/common"
)

type BundledBoilerplate struct{}

var _ TemplateRenderer = (*BundledBoilerplate)(nil)

func NewBundledRenderer() *BundledBoilerplate {
	return &BundledBoilerplate{}
}

func (b *BundledBoilerplate) Render(ctx context.Context, pkgManifestFilename string, outputFolders []string, baseUrlOrPath string) error {
	baseUrlOrPath = GetBaseUrlOrDefault(baseUrlOrPath)

	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return fmt.Errorf("loading package manifest: %w", err)
	}
	packagePathPrefix := GetPackagePathPrefixOrDefault(manifest.DefaultPackagePathPrefix)

	// Filter packages based on output folders
	var packagesToInstall []common.Package
	if len(outputFolders) == 0 {
		packagesToInstall = manifest.Packages
	} else {
		packagesToInstall = filterPackages(manifest.Packages, outputFolders)
	}

	for _, pkg := range packagesToInstall {
		templateURL := getValidTemplareUrlOrPath(baseUrlOrPath, pkg.Template, packagePathPrefix, pkg.Ref)
		slog.Debug("generating boilerplate template",
			slog.String("outputFolder", pkg.OutputFolder),
			slog.String("templateUrl", templateURL),
			slog.String("vars", strings.Join(pkg.VarFiles, ",")),
		)
		vars, err := variables.ParseVars(nil, pkg.VarFiles)
		if err != nil {
			return fmt.Errorf("parsing vars: %w", err)
		}
		opts := &options.BoilerplateOptions{
			TemplateUrl:    templateURL,
			OutputFolder:   pkg.OutputFolder,
			Vars:           vars,
			NonInteractive: true,
			//OnMissingKey:    options.DefaultMissingKeyAction,
			//OnMissingConfig: options.DefaultMissingConfigAction,
		}

		printPrettyOptionsAsCmd(opts, pkg.VarFiles)
		// root template does not have any dependencies
		emptyDependency := variables.Dependency{}
		if err := templates.ProcessTemplate(opts, opts, emptyDependency); err != nil {
			return fmt.Errorf("processing boilerplate template: %w", err)
		}
	}

	return nil
}

func printPrettyOptionsAsCmd(opts *options.BoilerplateOptions, varFiles []string) {
	cmdString := createPrettyCmdStringFromOpts(opts, varFiles)
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println("Running boilerplate command:")
	fmt.Println(green.Render(cmdString))
	fmt.Println("------------------------------------------------------------------------------------------")

}

func createPrettyCmdStringFromOpts(opts *options.BoilerplateOptions, varFiles []string) string {
	cmdString := "boilerplate"
	cmdString += fmt.Sprintf("\n  --template-url %s", opts.TemplateUrl)
	cmdString += fmt.Sprintf("\n  --output-folder %s", opts.OutputFolder)
	cmdString += fmt.Sprintf("\n  --non-interactive")
	for _, varFile := range varFiles {
		cmdString += fmt.Sprintf("\n  --var-file %s", varFile)
	}
	return cmdString
}
