package add

import (
	"fmt"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

type AddResult struct {
	OutputFolder    string
	VarFiles        []string
	TemplateName    string
	TemplateVersion string
}

/**
 * Add Boilerplate template to packages manifest with an optional stack name.
 * The template version is fetched from the latest release on GitHub and added to the packages manifest without applying the template.
 * The output folder is prefixed with the stack name and added to the packages manifest.
 */
func Run(pkgManifestFilename string, templateName, outputFolder string) (*AddResult, error) {

	latestReleases, err := githubreleases.GetLatestReleases()
	if err != nil {
		return nil, fmt.Errorf("failed getting latest github releases: %w", err)
	}

	templateVersion := latestReleases[templateName]
	if templateVersion == "" {
		return nil, fmt.Errorf("template %s not found in latest releases", templateName)
	}

	varFiles := []string{
		"_config/common-config.yml",
		fmt.Sprintf("_config/%s.yml", outputFolder),
	}

	newPackage := common.Package{
		Template:     templateName,
		Ref:          fmt.Sprintf("%s-%s", templateName, templateVersion),
		OutputFolder: outputFolder,
		VarFiles:     varFiles,
	}

	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return nil, err
	}

	for _, pkg := range manifest.Packages {
		if pkg.OutputFolder == newPackage.OutputFolder {
			return nil, fmt.Errorf("output folder %s already exists in manifest", newPackage.OutputFolder)
		}
	}

	manifest.Packages = append(manifest.Packages, newPackage)
	err = common.SavePackageManifest(pkgManifestFilename, manifest)
	if err != nil {
		return nil, err
	}

	return &AddResult{
		OutputFolder:    outputFolder,
		VarFiles:        varFiles,
		TemplateName:    templateName,
		TemplateVersion: templateVersion,
	}, nil
}
