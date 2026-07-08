package common

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ManifestExists returns true if a package manifest exists at the given path.
//
// LoadPackageManifest returns an empty manifest when the file does not exist, so callers that
// need to distinguish between updating an existing manifest and creating a new one should use
// this function.
func ManifestExists(filePath string) (bool, error) {
	return fileExists(filePath)
}

func LoadPackageManifest(filePath string) (PackageManifest, error) {
	fileContents, err := os.ReadFile(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return PackageManifest{}, nil
	}
	if err != nil {
		return PackageManifest{}, fmt.Errorf("opening manifest file '%s': %w", filePath, err)
	}

	var manifest PackageManifest

	err = yaml.Unmarshal(fileContents, &manifest)
	if err != nil {
		return PackageManifest{}, fmt.Errorf("unmarshalling YAML file '%s': %w", filePath, err)
	}

	return manifest, nil
}
