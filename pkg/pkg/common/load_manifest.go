package common

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadPackageManifest(filePath string) (PackageManifest, error) {
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return PackageManifest{}, fmt.Errorf("opening file: %w", err)
	}

	var manifest PackageManifest

	err = yaml.Unmarshal(fileContents, &manifest)
	if err != nil {
		return PackageManifest{}, fmt.Errorf("unmarshalling YAML: %w", err)
	}

	return manifest, nil
}
