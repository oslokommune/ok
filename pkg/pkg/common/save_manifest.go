package common

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func SavePackageManifest(filePath string, pm PackageManifest) (err error) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open manifest file '%s': %w", filePath, err)
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	defer enc.Close()
	if err := enc.Encode(pm); err != nil {
		return fmt.Errorf("failed to encode YAML manifest and write to file '%s': %w", filePath, err)
	}

	return nil
}
