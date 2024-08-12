package common

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadPackageManifest(filePath string) (PackageManifest, error) {
	fileContents, err := os.ReadFile(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return PackageManifest{}, nil
	}
	if err != nil {
		return PackageManifest{}, fmt.Errorf("error opening manifest file %s: %w", filePath, err)
	}

	var manifest PackageManifest

	err = yaml.Unmarshal(fileContents, &manifest)
	if err != nil {
		return PackageManifest{}, fmt.Errorf("error unmarshalling YAML file %s: %w", filePath, err)
	}

	return manifest, nil
}
