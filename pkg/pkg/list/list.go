package list

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
)

func Run(pkgManifestFilename string) ([]common.Package, error) {
	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return nil, fmt.Errorf("loading package manifest: %w", err)
	}

	return manifest.Packages, nil
}
