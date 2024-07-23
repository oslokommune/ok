package list

import (
	"fmt"

	"github.com/oslokommune/ok/pkg/pkg/update"
)

type Release struct {
	Component string
	Version   string
}

func Run(pkgManifestFilename string) error {
	manifest, err := update.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return fmt.Errorf("loading package manifest: %w", err)
	}

	fmt.Println(manifest)

	return nil
}
