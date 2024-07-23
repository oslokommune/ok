package list

import (
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
)

func Run(pkgManifestFilename string) error {
	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return fmt.Errorf("loading package manifest: %w", err)
	}

	fmt.Println(manifest)

	return nil
}
