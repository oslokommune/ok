package update

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"

	"gopkg.in/yaml.v3"

	"os"
)

type Release struct {
	Component string
	Version   string
}

func Run(pkgManifestFilename string) error {
	manifest, err := loadPackageManifest(pkgManifestFilename)
	if err != nil {
		return fmt.Errorf("loading package manifest: %w", err)
	}

	latestReleases, err := getLatestReleases()
	if err != nil {
		return fmt.Errorf("getting latest releases: %w", err)
	}

	// Set each package to the latest release
	for i, pkg := range manifest.Packages {
		manifest.Packages[i].Ref = fmt.Sprintf("%s-%s", pkg.Template, latestReleases[pkg.Template])
	}

	err = writePackageManifest(pkgManifestFilename, manifest)
	if err != nil {
		return fmt.Errorf("writing package manifest: %w", err)
	}

	return nil
}

func loadPackageManifest(filePath string) (common.PackageManifest, error) {
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return common.PackageManifest{}, fmt.Errorf("opening file: %w", err)
	}

	var manifest common.PackageManifest

	err = yaml.Unmarshal(fileContents, &manifest)
	if err != nil {
		return common.PackageManifest{}, fmt.Errorf("unmarshalling YAML: %w", err)
	}

	return manifest, nil
}

func writePackageManifest(pkgManifestFilename string, manifest common.PackageManifest) error {
	updatedYAML, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("error marshaling updated manifest: %w", err)
	}

	fileInfo, err := os.Stat(pkgManifestFilename)
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	err = os.WriteFile(pkgManifestFilename, updatedYAML, fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("error writing updated manifest to file: %w", err)
	}

	return nil
}
