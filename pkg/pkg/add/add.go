package add

import (
	"fmt"
	"path/filepath"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

/**
 * Add Boilerplate template to packages manifest with an optional stack name.
 * The template version is fetched from the latest release on GitHub and added to the packages manifest without applying the template.
 * The output folder is prefixed with the stack name and added to the packages manifest.
 */
func Run(pkgManifestFilename string, templateName, outputFolderName, stackName string) error {

	latestReleases, err := githubreleases.GetLatestReleases()
	if err != nil {
		return fmt.Errorf("failed getting latest github releases: %w", err)
	}

	templateVersion := latestReleases[templateName]
	if templateVersion == "" {
		return fmt.Errorf("template %s not found in latest releases", templateName)
	}

	outputFolder := filepath.Join(outputFolderName, stackName)
	varFiles := []string{
		"_config/common-config.yml",
		fmt.Sprintf("_config/%s.yml", stackName),
	}
	newPackage := common.Package{
		Template:     templateName,
		Ref:          fmt.Sprintf("%s-%s", templateName, templateVersion),
		OutputFolder: outputFolder,
		VarFiles:     varFiles,
	}

	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return err
	}

	for _, pkg := range manifest.Packages {
		if pkg.OutputFolder == newPackage.OutputFolder {
			return fmt.Errorf("output folder %s already exists in manifest", newPackage.OutputFolder)
		}
	}

	manifest.Packages = append(manifest.Packages, newPackage)
	return common.SavePackageManifest(pkgManifestFilename, manifest)
}
