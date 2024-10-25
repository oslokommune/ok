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

func SelectPackagesToInstall(pkgManifestFilename string) ([]string, error) {
	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return []string{}, fmt.Errorf("loading package manifest: %w", err)
	}

	options := make([]huh.Option[string], 0)

	for _, pkg := range manifest.Packages {
		displayText := createDisplayText(pkg)
		value := pkg.OutputFolder

		options = append(options, huh.NewOption[string](displayText, value))
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

func createDisplayText(pkg common.Package) string {
	/*
	   Format specifiers explained:
	   %-*.*s formats a string with a specified width and precision.
	   - The first * sets the minimum width, padding with spaces if shorter.
	   - The second * sets the maximum length, truncating if longer.
	*/
	outputFolder := fmt.Sprintf("%-*.*s", outputFolderWidth, outputFolderWidth, pkg.OutputFolder)
	template := fmt.Sprintf("%-*.*s", templateWidth, templateWidth, pkg.Template)
	varFiles := fmt.Sprintf("%-*.*s", varFilesWidth, varFilesWidth, fmt.Sprint(pkg.VarFiles))

	return fmt.Sprintf("%s %s %s", outputFolder, template, varFiles)
}
