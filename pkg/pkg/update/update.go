package update

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/oslokommune/ok/pkg/pkg/schema"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"
	"os"
	"path"
	"strings"

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

		varFile, existingSchema, found := findSchemaVarFile(pkg, workingDirectory)
		if !found {
			_, _ = fmt.Fprintf(os.Stderr,
				"⚠️ Warning: no var file of package '%s' declares a JSON schema for template '%s'."+
					" Skipping schema update for this package."+
					" Add a '# yaml-language-server: $schema=...' declaration to the package's own var file to enable schema updates.\n",
				pkg.OutputFolder, pkg.Template)
			continue
		}

		fmt.Printf("- %s\n", pkg.OutputFolder)

		if existingSchema.Ref() == pkg.Ref {
			// No need to update the varFile with a new JSON schema, as the existing one is as declared in the package
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

// findSchemaVarFile returns the var file that should receive the package's JSON schema declaration.
//
// The file is selected by looking at the schema declaration each var file already carries, instead of relying on the
// position of the file in pkg.VarFiles (the position is meaningful to Boilerplate, but says nothing about which file
// is the package's own config). The rules are:
//
//   - A var file whose declaration references the package's template is a candidate. If several var files qualify,
//     the last one wins.
//   - A var file that declares a different template (for example a shared common-config.yml declared by another
//     package) is never selected, as overwriting its declaration would break schema validation for the other users
//     of the file.
//   - A var file without a (parseable) schema declaration is never selected. `ok pkg add` scaffolds the declaration
//     into the package's var file, so guessing by position here would risk writing a package-specific declaration
//     into a shared file.
//
// If no var file qualifies, found is false.
func findSchemaVarFile(pkg common.Package, workingDirectory string) (varFile string, existingSchema metadata.JsonSchema, found bool) {
	for _, varFileRelative := range pkg.VarFiles {
		candidate := path.Join(workingDirectory, varFileRelative)

		jsonSchema, err := metadata.ParseFirstLine(candidate)
		if err != nil {
			// The file is missing, empty, or has no parseable schema declaration on its first line, so we cannot
			// tell if it belongs to this package.
			continue
		}

		if jsonSchema.Template != pkg.Template {
			// The file declares a schema for a different template, e.g. a shared config file.
			continue
		}

		varFile = candidate
		existingSchema = jsonSchema
		found = true
	}

	return varFile, existingSchema, found
}

type Options struct {
	DisableManifestUpdate bool
	MigrateConfig         bool
	UpdateSchema          bool
}
