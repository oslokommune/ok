package interactive

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/oslokommune/ok/pkg/pkg/common"
)

func SelectPackages(manifest common.PackageManifest, action string) ([]common.Package, error) {
	options := make([]huh.Option[string], 0)
	packageMap := make(map[string]common.Package)

	for _, pkg := range manifest.Packages {
		displayText := pkg.String()

		options = append(options, huh.NewOption[string](displayText, pkg.Key()))
		packageMap[pkg.Key()] = pkg
	}

	var selectedPackageKeys []string

	s := huh.NewMultiSelect[string]().
		Options(options...).
		Title("Select package(s) to " + action).
		Limit(4).
		Value(&selectedPackageKeys)

	err := huh.NewForm(huh.NewGroup(s)).Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return []common.Package{}, nil
		} else {
			return []common.Package{}, fmt.Errorf("running multi select form: %w", err)
		}
	}

	var selectedPackages []common.Package
	for _, key := range selectedPackageKeys {
		pkg := packageMap[key]
		selectedPackages = append(selectedPackages, pkg)
	}

	return selectedPackages, nil
}
