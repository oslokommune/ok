package update

import (
	"fmt"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

func Run(pkgManifestFilename string) error {
	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return fmt.Errorf("loading package manifest: %w", err)
	}

	latestReleases, err := githubreleases.GetLatestReleases()
	if err != nil {
		return fmt.Errorf("getting latest releases: %w", err)
	}

	// Set each package to the latest release
	for i, pkg := range manifest.Packages {
		manifest.Packages[i].Ref = fmt.Sprintf("%s-%s", pkg.Template, latestReleases[pkg.Template])
	}

	err = common.SavePackageManifest(pkgManifestFilename, manifest)
	if err != nil {
		return err
	}

	return nil
}
