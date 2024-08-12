package add

import (
	"fmt"
	"path/filepath"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

func Run(pkgManifestFilename string, templateName, outputFolderName, appName string) error {

	latestReleases, err := githubreleases.GetLatestReleases()
	if err != nil {
		return fmt.Errorf("failed getting latest github releases: %w", err)
	}

	templateVersion := latestReleases[templateName]
	if templateVersion == "" {
		return fmt.Errorf("template %s not found in latest releases", templateName)
	}

	templateAppName := templateName
	if appName != "" {
		templateAppName = fmt.Sprintf("%s-%s", templateName, appName)
	}
	outputFolder := filepath.Join(outputFolderName, templateAppName)
	varFiles := []string{
		"config/common-config.yml",
		fmt.Sprintf("config/%s.yml", templateAppName),
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
