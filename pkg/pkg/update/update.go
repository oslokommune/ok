package update

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/ok/pkg/pkg/schema"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"

	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/githubreleases"
)

type Updater struct {
	ghReleases GitHubReleases
}

type GitHubReleases interface {
	GetLatestReleases() (map[string]string, error)
}

func NewUpdater(ghReleases GitHubReleases) Updater {
	return Updater{
		ghReleases: ghReleases,
	}
}

// Run updates the package manifest with the latest releases.
func (u Updater) Run(pkgManifestFilename string, selectedPackagesInput []common.Package, opts Options, workingDirectory string) error {
	var selectedPackages []common.Package
	var manifest common.PackageManifest
	{
		currentManifest, err := common.LoadPackageManifest(pkgManifestFilename)
		if err != nil {
			return fmt.Errorf("loading package manifest: %w", err)
		}

		if opts.DisableManifestUpdate {
			manifest = currentManifest
			selectedPackages = selectedPackagesInput
		} else {
			updateManifest, updatedSelectedPackages, err := u.updateManifest(currentManifest, selectedPackagesInput)
			if err != nil {
				return fmt.Errorf("updating package manifest: %w", err)
			}

			manifest = updateManifest
			selectedPackages = updatedSelectedPackages

			err = common.SavePackageManifest(pkgManifestFilename, manifest)
			if err != nil {
				return fmt.Errorf("saving package manifest: %w", err)
			}
		}
	}

	if opts.UpdateSchema {
		err := u.setJsonSchemaDeclarationInVarFiles(selectedPackages, workingDirectory)
		if err != nil {
			return err
		}
	}

	if opts.MigrateConfig {
		err := migrate_config.MigrateVarFile(selectedPackages, workingDirectory)
		if err != nil {
			return fmt.Errorf("migrating package config: %w", err)
		}
	} else {
		fmt.Println("Not migrating package configuration files.")
	}

	common.PrintProcessedPackages(selectedPackages, "updated")

	return nil
}

// updateManifest updates package versions in the manifest with the latest releases from GitHub
func (u Updater) updateManifest(manifest common.PackageManifest, selectedPackages []common.Package) (common.PackageManifest, []common.Package, error) {
	latestReleases, err := u.ghReleases.GetLatestReleases()
	if err != nil {
		if strings.Contains(err.Error(), "secret not found in keyring") {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", githubreleases.AuthErrorHelpMessage)
		}

		return common.PackageManifest{}, nil, fmt.Errorf("failed getting latest github releases: %w", err)
	}

	updatedManifest, updatedPackages, err := updatePackages(manifest, selectedPackages, latestReleases)
	if err != nil {
		return common.PackageManifest{}, nil, fmt.Errorf("updating packages: %w", err)
	}

	return updatedManifest, updatedPackages, nil
}

// updatePackages updates the package manifest with the latest releases. It only updates the packages found in selectedPackages.
func updatePackages(manifest common.PackageManifest, selectedPackages []common.Package, latestReleases map[string]string) (common.PackageManifest, []common.Package, error) {
	updatedManifest := manifest.Clone()
	updatedPackages := make([]common.Package, 0)

	for i, _ := range manifest.Packages {
		pkg := &updatedManifest.Packages[i]

		_, err := pkg.PackageVersion()
		if errors.Is(err, semver.ErrInvalidSemVer) {
			// pkg.Ref is not a semver version, e.g. "main". We don't update these.
			continue
		}

		if !common.ContainsPackage(selectedPackages, *pkg) {
			continue
		}

		// pkg is a package that is in selectedPackages, i.e. it should be updated
		latestRelease, ok := latestReleases[pkg.Template] // e.g. v2.1.3
		if !ok {
			return common.PackageManifest{}, nil, fmt.Errorf("no latest release found for package: %s", pkg.Template)
		}

		newRef := fmt.Sprintf("%s-%s", pkg.Template, latestRelease) // e.g. app-v2.1.3
		pkg.Ref = newRef
		updatedPackages = append(updatedPackages, *pkg)
	}

	return updatedManifest, updatedPackages, nil
}

// setJsonSchemaDeclarationInVarFiles sets the varFile's JSON schema declaration to the same version as defined in the
// package manifest.
func (u Updater) setJsonSchemaDeclarationInVarFiles(selectedPackages []common.Package, workingDirectory string) error {
	fmt.Println("Updating json schemas:")

	for _, pkg := range selectedPackages {
		_, err := pkg.PackageVersion()
		if err != nil {
			// pkg.Ref might be "main". To keep code simple, we avoid dealing with non-semver versions.
			continue
		}

		varFile, ok := getLastVarFile(pkg, workingDirectory)
		if !ok {
			continue
		}

		fmt.Printf("- %s\n", pkg.OutputFolder)

		// Get current JSON schema for the package
		jsonSchemaMetdata, err := metadata.ParseFirstLine(varFile)
		if err != nil && errors.Is(err, metadata.ErrMissingSchemaDeclaration) {
			err = schema.SetSchemaDeclarationInVarFile(varFile, pkg.Ref)
			if err != nil {
				return fmt.Errorf("creating or updating configuration file: %w", err)
			}

			return nil
		} else if err != nil {
			return fmt.Errorf("parsing first line of file '%s': %w", varFile, err)
		}

		existingRef := fmt.Sprintf("%s-v%s", jsonSchemaMetdata.Template, jsonSchemaMetdata.Version)
		if existingRef == pkg.Ref {
			// No need to update the varFile with a new JSON schema, as the existing one is as declared in the pacckage
			// manifest.
			continue
		}

		err = schema.SetSchemaDeclarationInVarFile(varFile, pkg.Ref)
		if err != nil {
			return fmt.Errorf("creating or updating configuration file: %w", err)
		}
	}

	return nil
}

func getLastVarFile(pkg common.Package, workingDirectory string) (string, bool) {
	if len(pkg.VarFiles) > 0 {
		varFileRelative := pkg.VarFiles[len(pkg.VarFiles)-1]
		varFile := path.Join(workingDirectory, varFileRelative)
		return varFile, true
	}

	return "", false
}

type Options struct {
	DisableManifestUpdate bool
	MigrateConfig         bool
	UpdateSchema          bool
}
