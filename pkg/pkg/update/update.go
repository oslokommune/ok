package update

import (
	"fmt"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

func Run(pkgManifestFilename string, packageName string) error {
	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return fmt.Errorf("loading package manifest: %w", err)
	}

	latestReleases, err := githubreleases.GetLatestReleases()
	if err != nil {
		return fmt.Errorf("getting latest releases: %w", err)
	}

	if packageName != "" {
		// Update only the specified package
		updated := false
		for i, pkg := range manifest.Packages {
			if pkg.Template == packageName {
				latestRelease, ok := latestReleases[pkg.Template]
				if !ok {
					return fmt.Errorf("no latest release found for package: %s", packageName)
				}
				manifest.Packages[i].Ref = fmt.Sprintf("%s-%s", pkg.Template, latestRelease)
				updated = true
				break
			}
		}
		if !updated {
			return fmt.Errorf("package not found: %s", packageName)
		}
	} else {
		// Update all packages
		for i, pkg := range manifest.Packages {
			latestRelease, ok := latestReleases[pkg.Template]
			if !ok {
				return fmt.Errorf("no latest release found for package: %s", pkg.Template)
			}
			manifest.Packages[i].Ref = fmt.Sprintf("%s-%s", pkg.Template, latestRelease)
		}
	}

	err = common.SavePackageManifest(pkgManifestFilename, manifest)
	if err != nil {
		return fmt.Errorf("saving package manifest: %w", err)
	}

	return nil
}
