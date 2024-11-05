package interactive

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/oslokommune/ok/pkg/pkg/common"
)

const (
	outputFolderWidth = 45
	templateWidth     = 40
	varFilesWidth     = 80
)

func SelectPackagesToInstall(manifest common.PackageManifest) ([]common.Package, error) {
	options := make([]huh.Option[string], 0)
	packageMap := make(map[string]common.Package)

	for _, pkg := range manifest.Packages {
		displayText := createDisplayText(pkg)

		options = append(options, huh.NewOption[string](displayText, pkg.Key()))
		packageMap[pkg.Key()] = pkg
	}

	var selectedPackageKeys []string

	s := huh.NewMultiSelect[string]().
		Options(options...).
		Title("Select package(s) to install").
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

func createDisplayText(pkg common.Package) string {
	outputFolder := fmt.Sprintf("%-*.*s", outputFolderWidth, outputFolderWidth, pkg.OutputFolder)
	template := fmt.Sprintf("%-*.*s", templateWidth, templateWidth, pkg.Template)
	varFiles := fmt.Sprintf("%-*.*s", varFilesWidth, varFilesWidth, fmt.Sprint(pkg.VarFiles))

	return fmt.Sprintf("%s %s %s", outputFolder, template, varFiles)
}
