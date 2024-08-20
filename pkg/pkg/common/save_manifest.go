package common

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func SavePackageManifest(filePath string, manifest PackageManifest) (err error) {
	// if the file already exists it will be overwritten
	// and the file PERM will be retained
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("opening manifest file '%s': %w", filePath, err)
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	defer enc.Close()
	if err := enc.Encode(manifest); err != nil {
		return fmt.Errorf("encoding YAML manifest and write to file '%s': %w", filePath, err)
	}

	return nil
}
