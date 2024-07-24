package interactive

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/oslokommune/ok/pkg/pkg/common"
)

func SelectPackagesToInstall(pkgManifestFilename string) ([]string, error) {
	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return []string{}, fmt.Errorf("loading package manifest: %w", err)
	}

	options := make([]huh.Option[string], 0)

	for _, p := range manifest.Packages {
		options = append(options, huh.NewOption(p.OutputFolder, p.OutputFolder))
	}

	var packages []string

	s := huh.NewMultiSelect[string]().
		Options(options...).
		Title("Select package(s) to install").
		Limit(4).
		Value(&packages)

	err = huh.NewForm(huh.NewGroup(s)).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return []string{}, nil
		} else {
			return []string{}, fmt.Errorf("running multi select form: %w", err)
		}
	}

	return packages, nil
}
