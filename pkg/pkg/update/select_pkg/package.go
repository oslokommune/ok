package select_pkg

import (
	"fmt"
	"os"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"gopkg.in/yaml.v3"
)

func listPackages() ([]string, error) {
	manifest, err := loadPackageManifest("packages.yml")
	if err != nil {
		return nil, fmt.Errorf("loading package manifest: %w", err)
	}

	var packages []string
	for _, pkg := range manifest.Packages {
		packages = append(packages, pkg.Template)
	}

	return packages, nil
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
